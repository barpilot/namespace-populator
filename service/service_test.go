package service

import (
	"fmt"
	"sync"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type logKind int

const (
	infoKind logKind = iota
	warnignKind
	errorKind
)

type logEvent struct {
	kind logKind
	line string
}

type testLogger struct {
	events []logEvent
	sync.Mutex
}

func (t *testLogger) logLine(kind logKind, format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	t.events = append(t.events, logEvent{kind: kind, line: str})
}

func (t *testLogger) Infof(format string, args ...interface{}) {
	t.logLine(infoKind, format, args...)
}
func (t *testLogger) Warningf(format string, args ...interface{}) {
	t.logLine(warnignKind, format, args...)
}
func (t *testLogger) Errorf(format string, args ...interface{}) {
	t.logLine(errorKind, format, args...)
}

// func TestEchoServiceEchoString(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		prefix     string
// 		msg        string
// 		expResults []logEvent
// 	}{
// 		{
// 			name:   "Logging a prefix and a string should log.",
// 			prefix: "test",
// 			msg:    "this is a test",
// 			expResults: []logEvent{
// 				logEvent{kind: infoKind, line: "[test] this is a test"},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			assert := assert.New(t)

// 			// Mocks.
// 			ml := &testLogger{events: []logEvent{}}

// 			// Create aservice and run.
// 			srv := service.NewSimpleEcho(ml)
// 			srv.EchoS(test.prefix, test.msg)

// 			// Check.
// 			assert.Equal(test.expResults, ml.events)
// 		})
// 	}
// }

// func TestEchoServiceEchoObj(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		prefix     string
// 		obj        runtime.Object
// 		expResults []logEvent
// 	}{
// 		{
// 			name:   "Logging a pod should print pod name.",
// 			prefix: "test",
// 			obj: &corev1.Pod{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name: "mypod",
// 				},
// 			},
// 			expResults: []logEvent{
// 				logEvent{kind: infoKind, line: "[test] mypod"},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			assert := assert.New(t)

// 			// Mocks.
// 			ml := &testLogger{events: []logEvent{}}

// 			// Create aservice and run.
// 			srv := service.NewSimpleEcho(ml)
// 			srv.EchoObj(test.prefix, test.obj)

// 			// Check.
// 			assert.Equal(test.expResults, ml.events)
// 		})
// 	}
// }

func Test_validateConfigMap(t *testing.T) {
	type args struct {
		cm *corev1.ConfigMap
		ns *corev1.Namespace
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "without selector",
			args: args{
				cm: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mycm",
					},
				},
				ns: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "myns",
					},
				},
			},
			want: true,
		},
		{
			name: "with selector",
			args: args{
				cm: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mycm",
						Annotations: map[string]string{
							"namespace-populator.barpilot.io/selector": "app=test",
						},
					},
				},
				ns: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "myns",
						Labels: map[string]string{
							"app": "test",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "with incorrect selector",
			args: args{
				cm: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mycm",
						Annotations: map[string]string{
							"namespace-populator.barpilot.io/selector": "app=test",
						},
					},
				},
				ns: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "myns",
						Labels: map[string]string{
							"app": "notgood",
						},
					},
				},
			},
			want: false,
		},
		{
			name: "with bad selector",
			args: args{
				cm: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mycm",
						Annotations: map[string]string{
							"namespace-populator.barpilot.io/selector": "apptest",
						},
					},
				},
				ns: &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: "myns",
						Labels: map[string]string{
							"app": "test",
						},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateConfigMap(tt.args.cm, tt.args.ns); got != tt.want {
				t.Errorf("validateConfigMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
