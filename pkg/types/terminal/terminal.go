package terminal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"time"
)

var Terminal terminal

type terminal struct{}

var upgrader = func() websocket.Upgrader {
	u := websocket.Upgrader{}
	u.HandshakeTimeout = time.Second * 2
	u.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return u
}()

func (t *TerminalSession) Done() {
	close(t.DoneChan)
}

func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.SizeChan:
		return &size
	case <-t.DoneChan:
		return nil
	}
}

func (t *TerminalSession) Read(p []byte) (int, error) {
	_, data, err := t.WsConn.ReadMessage()
	if err != nil {
		return 0, err
	}
	var msg TerminalMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return 0, err
	}
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.SizeChan <- remotecommand.TerminalSize{
			Width:  msg.Cols,
			Height: msg.Rows,
		}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		return 0, errors.New(fmt.Sprintf("unknown operation type: %s\n", msg.Operation))
	}
}

func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		return 0, errors.New(fmt.Sprintf("write parse message err: %v\n", err))
	}

	if err := t.WsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, errors.New(fmt.Sprintf("write message err: %v\n", err))
	}
	return len(p), nil
}

func (t *TerminalSession) Close() error {
	return t.WsConn.Close()
}

func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	cluster := r.Form.Get("cluster")
	fmt.Printf("exec pod: %s, container: %s, namespace: %s, cluster: %s", podName, containerName, namespace, cluster)

	client, err := utils.K8s.GetClient(cluster)
	if err != nil {
		return
	}

	conf, err := clientcmd.BuildConfigFromFlags("", utils.K8s.KubeConfMap[cluster])
	if err != nil {
		return
	}

	pty, err := NewTerminalSession(w, r, nil)
	if err != nil {
		return
	}

	defer func() {
		err := pty.Close()
		if err != nil {
			return
		}
	}()

	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		return
	}

	if err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	}); err != nil {
		pty.Write([]byte(fmt.Sprintf("Exec to pod error: %v\n", err)))
		pty.Done()
	}

}

func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	wsConn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	return &TerminalSession{
		WsConn:   wsConn,
		SizeChan: make(chan remotecommand.TerminalSize),
		DoneChan: make(chan struct{}),
	}, nil
}
