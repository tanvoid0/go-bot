package commands

import (
	"context"
	"log"

	"github.com/rivo/tview"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Resource struct {
	Type string
	Name string
}

var selectedResources = make(map[string]bool)

func RunK10s() {
	// Initialize Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Application and UI setup
	app := tview.NewApplication()

	// Drop-down for selecting namespace
	namespaceDropDown := tview.NewDropDown().
		SetLabel("Select Namespace: ").
		SetOptions([]string{"default", "kube-system"}, nil)

	// Table for showing pods and statefulsets with checkboxes
	table := tview.NewTable().
		SetSelectable(true, false).
		SetBorders(true)

	// Fetch pods and statefulsets from selected namespace and populate table
	fetchResources := func(namespace string) {
		// Clear the table
		table.Clear()

		// Set table headers
		table.SetCell(0, 0, tview.NewTableCell("Select").SetSelectable(false))
		table.SetCell(0, 1, tview.NewTableCell("Type").SetSelectable(false))
		table.SetCell(0, 2, tview.NewTableCell("Name").SetSelectable(false))

		// Fetch pods
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error fetching pods: %v", err)
			return
		}

		// Fetch statefulsets
		statefulsets, err := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Printf("Error fetching statefulsets: %v", err)
			return
		}

		row := 1
		// Add pods to the table
		for _, pod := range pods.Items {
			table.SetCell(row, 0, tview.NewTableCell("[ ]").SetReference(pod.Name))
			table.SetCell(row, 1, tview.NewTableCell("Pod"))
			table.SetCell(row, 2, tview.NewTableCell(pod.Name))
			row++
		}

		// Add statefulsets to the table
		for _, ss := range statefulsets.Items {
			table.SetCell(row, 0, tview.NewTableCell("[ ]").SetReference(ss.Name))
			table.SetCell(row, 1, tview.NewTableCell("StatefulSet"))
			table.SetCell(row, 2, tview.NewTableCell(ss.Name))
			row++
		}
	}

	// On namespace selection, fetch and display resources
	namespaceDropDown.SetSelectedFunc(func(text string, index int) {
		fetchResources(text)
	})

	// Set behavior for selecting/deselecting resources via checkboxes
	table.SetSelectedFunc(func(row, column int) {
		if row > 0 {
			cell := table.GetCell(row, 0)
			name := table.GetCell(row, 2).Text

			// Toggle checkbox
			if cell.Text == "[ ]" {
				cell.SetText("[X]")
				selectedResources[name] = true
			} else {
				cell.SetText("[ ]")
				delete(selectedResources, name)
			}
		}
	})

	// Layout: Dropdown + Table
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(namespaceDropDown, 1, 0, true).
		AddItem(table, 0, 1, true)

	// Quit application with ESC
	app.SetRoot(layout, true).SetFocus(namespaceDropDown).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				app.Stop()
			}
			return event
		})

	// Run application
	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
