// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event

import (
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api"

	"github.com/kubernetes/dashboard/resource/common"
)

func TestGetPodsEventWarningsApi(t *testing.T) {
	cases := []struct {
		pods      []api.Pod
		eventList []api.Event
		expected  []common.Event
	}{
		{nil, nil, []common.Event{}},
		{
			[]api.Pod{
				{
					ObjectMeta: api.ObjectMeta{
						Name: "FailedPod",
					},
					Status: api.PodStatus{
						Phase: api.PodFailed,
					},
				},
			},
			[]api.Event{
				{
					Type:    api.EventTypeWarning,
					Message: "Test Message",
					InvolvedObject: api.ObjectReference{
						Name: "FailedPod",
					},
				},
			},
			[]common.Event{
				{
					Message: "Test Message",
					Type:    api.EventTypeWarning,
				},
			},
		},
		{
			[]api.Pod{
				{
					Status: api.PodStatus{
						Phase: api.PodRunning,
					},
				},
			},
			nil,
			[]common.Event{},
		},
	}

	for _, c := range cases {
		actual := GetPodsEventWarnings(c.eventList, c.pods)

		if len(actual) != len(c.expected) || !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("GetPodsEventWarnings(%#v, %#v) == \n%#v\nexpected \n%#v\n",
				c.eventList, c.pods, actual, c.expected)
		}
	}
}

func TestGetWarningEvents(t *testing.T) {
	cases := []struct {
		events   *api.EventList
		expected []api.Event
	}{
		{&api.EventList{Items: []api.Event{}}, []api.Event{}},
		{
			&api.EventList{
				Items: []api.Event{
					{
						Message: "msg",
						Reason:  "reason",
						Type:    api.EventTypeWarning,
					},
				},
			},
			[]api.Event{
				{
					Message: "msg",
					Reason:  "reason",
					Type:    api.EventTypeWarning,
				},
			},
		},
		{
			&api.EventList{
				Items: []api.Event{
					{
						Message: "msg",
						Reason:  "failed",
					},
				},
			},
			[]api.Event{
				{
					Message: "msg",
					Reason:  "failed",
					Type:    api.EventTypeWarning,
				},
			},
		},
		{
			&api.EventList{
				Items: []api.Event{
					{
						Message: "msg",
						Reason:  "reason",
					},
				},
			},
			[]api.Event{},
		},
	}

	for _, c := range cases {
		actual := getWarningEvents(c.events.Items)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("getWarningEvents(%#v) == \n%#v\nexpected \n%#v\n",
				c.events, actual, c.expected)
		}
	}
}

func TestFilterEventsByType(t *testing.T) {
	events := []api.Event{
		{Type: api.EventTypeNormal},
		{Type: api.EventTypeWarning},
	}

	cases := []struct {
		events    []api.Event
		eventType string
		expected  []api.Event
	}{
		{nil, "", nil},
		{nil, api.EventTypeWarning, nil},
		{
			events,
			"",
			events,
		},
		{
			events,
			api.EventTypeNormal,
			[]api.Event{
				{Type: api.EventTypeNormal},
			},
		},
		{
			events,
			api.EventTypeWarning,
			[]api.Event{
				{Type: api.EventTypeWarning},
			},
		},
	}

	for _, c := range cases {
		actual := filterEventsByType(c.events, c.eventType)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("FilterEventsByType(%#v, %#v) == \n%#v\nexpected \n%#v\n",
				c.events, c.eventType, actual, c.expected)
		}
	}
}

func TestRemoveDuplicates(t *testing.T) {
	cases := []struct {
		slice    []api.Event
		expected []api.Event
	}{
		{nil, []api.Event{}},
		{
			[]api.Event{
				{Reason: "test"},
				{Reason: "test2"},
				{Reason: "test"},
			},
			[]api.Event{
				{Reason: "test"},
				{Reason: "test2"},
			},
		},
		{
			[]api.Event{
				{Reason: "test"},
				{Reason: "test"},
				{Reason: "test"},
			},
			[]api.Event{
				{Reason: "test"},
			},
		},
		{
			[]api.Event{
				{Reason: "test"},
				{Reason: "test2"},
				{Reason: "test3"},
			},
			[]api.Event{
				{Reason: "test"},
				{Reason: "test2"},
				{Reason: "test3"},
			},
		},
	}

	for _, c := range cases {
		actual := removeDuplicates(c.slice)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("removeDuplicates(%#v) == \n%#v\nexpected \n%#v\n",
				c.slice, actual, c.expected)
		}
	}
}

func TestIsRunningOrSucceeded(t *testing.T) {
	cases := []struct {
		pod      api.Pod
		expected bool
	}{
		{
			api.Pod{
				Status: api.PodStatus{
					Phase: api.PodRunning,
				},
			},
			true,
		},
		{
			api.Pod{
				Status: api.PodStatus{
					Phase: api.PodSucceeded,
				},
			},
			true,
		},
		{
			api.Pod{
				Status: api.PodStatus{
					Phase: api.PodFailed,
				},
			},
			false,
		},
		{
			api.Pod{
				Status: api.PodStatus{
					Phase: api.PodPending,
				},
			},
			false,
		},
	}

	for _, c := range cases {
		actual := isRunningOrSucceeded(c.pod)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("isRunningOrSucceded(%#v) == \n%#v\nexpected \n%#v\n",
				c.pod, actual, c.expected)
		}
	}
}

func TestIsTypeFilled(t *testing.T) {
	cases := []struct {
		events   []api.Event
		expected bool
	}{
		{nil, false},
		{
			[]api.Event{
				{Type: api.EventTypeWarning},
			},
			true,
		},
		{
			[]api.Event{},
			false,
		},
		{
			[]api.Event{
				{Type: api.EventTypeWarning},
				{Type: api.EventTypeNormal},
				{Type: ""},
			},
			false,
		},
	}

	for _, c := range cases {
		actual := IsTypeFilled(c.events)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("IsTypeFilled(%#v) == \n%#v\nexpected \n%#v\n",
				c.events, actual, c.expected)
		}
	}
}

func TestFillEventsType(t *testing.T) {
	cases := []struct {
		events   []api.Event
		expected []api.Event
	}{
		{nil, nil},
		{[]api.Event{}, []api.Event{}},
		{
			[]api.Event{
				{Reason: "failed"},
				{Reason: "test"},
			},
			[]api.Event{
				{Reason: "failed", Type: api.EventTypeWarning},
				{Reason: "test", Type: api.EventTypeNormal},
			},
		},
	}

	for _, c := range cases {
		actual := FillEventsType(c.events)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("FillEventsType(%#v) == \n%#v\nexpected \n%#v\n",
				c.events, actual, c.expected)
		}
	}
}

func TestIsFailedReason(t *testing.T) {
	cases := []struct {
		reason         string
		failedPartials []string
		expected       bool
	}{
		{"", nil, false},
		{"", []string{}, false},
		{"FailedReason", []string{"failed"}, true},
		{"ErrReason", []string{"failed"}, false},
		{"ErrReason", []string{"failed", "err"}, true},
	}

	for _, c := range cases {
		actual := IsFailedReason(c.reason, c.failedPartials...)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("IsFailedReason(%#v, %#v) == \n%#v\nexpected \n%#v\n",
				c.reason, c.failedPartials, actual, c.expected)
		}
	}
}

func TestFilterEventsByPodsUID(t *testing.T) {
	cases := []struct {
		events   []api.Event
		pods     []api.Pod
		expected []api.Event
	}{
		{nil, nil, []api.Event{}},
		{
			[]api.Event{
				{InvolvedObject: api.ObjectReference{UID: "TestPod"}},
				{InvolvedObject: api.ObjectReference{UID: "TestObject"}},
			},
			[]api.Pod{
				{ObjectMeta: api.ObjectMeta{UID: "TestPod"}},
			},
			[]api.Event{
				{InvolvedObject: api.ObjectReference{UID: "TestPod"}},
			},
		},
		{
			// To check whether multiple events targeting same object are correctly filtered
			[]api.Event{
				{InvolvedObject: api.ObjectReference{UID: "TestPod"}},
				{InvolvedObject: api.ObjectReference{UID: "TestPod"}},
			},
			[]api.Pod{
				{ObjectMeta: api.ObjectMeta{UID: "TestPod"}},
			},
			[]api.Event{
				{InvolvedObject: api.ObjectReference{UID: "TestPod"}},
				{InvolvedObject: api.ObjectReference{UID: "TestPod"}},
			},
		},
	}

	for _, c := range cases {
		actual := FilterEventsByPodsUID(c.events, c.pods)
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("FilterEventsByPodsUID(%#v, %#v) == \n%#v\nexpected \n%#v\n",
				c.events, c.pods, actual, c.expected)
		}
	}
}