/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/subhamBharadwaz/go-todo-cli-app/utils"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Marks a todo item as completed by ID",
	Long: `The complete command allows you to mark a todo item as completed in the CSV file 'tasks.csv'. You need to provide the ID of the task you want to mark as completed. This command updates the task's status to indicate that it has been completed.

Example usage:
  $ tasks complete 3

In this example, the task with ID 3 will be updated to reflect that it has been completed. If the task with the specified ID does not exist, an error message will be displayed. The task status will be updated in the 'tasks.csv' file, and the file will be modified to include the new status.

Note: Ensure that the task ID provided is valid. This command will not prompt for confirmation before marking the task as completed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide the ID of the task to complete.")
			return
		}

		id := args[0]
		task, err := completeTask(id)
		if err != nil {
			fmt.Println("Error completing task:", err)
		} else {
			writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
			fmt.Fprintln(writer, "ID\tName\tCreated\tDone\t")

			createdTime, _ := time.Parse(time.RFC3339, task[2])

			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t\n", task[0], task[1], timediff.TimeDiff(createdTime), task[3])

			writer.Flush()
		}
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}

func completeTask(id string) ([]string, error) {
	file, err := utils.LoadFile("tasks.csv")
	if err != nil {
		return nil, err
	}
	defer utils.CloseFile(file)

	// Read existing tasks
	tasks, err := utils.ReadTasks(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks: %w", err)
	}

	// Find task and add complete as true with the given ID
	var updatedTasks [][]string
	var updatedTask []string
	var found bool
	for _, task := range tasks {
		if task[0] == id {
			found = true
			if len(task) >= 4 {
				task[3] = "true"
			} else {
				return nil, fmt.Errorf("task format is invalid, no done column found")
			}
			updatedTask = task
		}
		updatedTasks = append(updatedTasks, task)
	}

	if !found {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}

	// Write updated task back to file
	if err := file.Truncate(0); err != nil {
		return nil, fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(updatedTasks); err != nil {
		return nil, fmt.Errorf("failed to write tasks: %w", err)
	}
	writer.Flush()
	return updatedTask, nil
}
