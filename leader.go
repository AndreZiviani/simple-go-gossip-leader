package main

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func leader(ctx context.Context, me string) *leaderelection.LeaderElector {
	// create the config from the config file
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Failed to create config: %v", err)
	}

	// create the clientset
	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create clientset: %v", err)
	}

	// create the lock config
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "leader-election-sample",
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: me,
		},
	}

	// create the leader election config
	leaderelectionConfig := leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: 5 * time.Second,
		RenewDeadline: 4 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				fmt.Println("I am the leader")
			},
			OnStoppedLeading: func() {
				fmt.Println("I am not the leader")
			},
			OnNewLeader: func(identity string) {
				fmt.Printf("New leader elected: %s\n", identity)
			},
		},
		ReleaseOnCancel: true,
	}

	// run the leader election
	le, error := leaderelection.NewLeaderElector(leaderelectionConfig)
	if error != nil {
		fmt.Printf("Failed to create leader elector: %v\n", error)
	}

	go le.Run(ctx)

	return le
}
