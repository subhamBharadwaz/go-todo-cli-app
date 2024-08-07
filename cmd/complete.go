/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
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

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID:", err)
			return
		}

		task, err := completeTask(id)
		if err != nil {
			fmt.Println("Error marking task as completed:", err)
			return
		}

		// Display the updated task in a tabular format
		writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
		fmt.Fprintln(writer, "ID\tName\tCreated\tDone\tDue\t")

		createdTime, _ := time.Parse(time.RFC3339, task.CreatedAt)
		completed := strconv.FormatBool(task.Completed)

		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\t\n", task.ID, task.Description, timediff.TimeDiff(createdTime), completed, task.DueDate)

		writer.Flush()

	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}

func completeTask(id int) (Task, error) {

	_, err := utils.DB.Exec("UPDATE tasks SET completed = ? WHERE id= ?", true, id)
	if err != nil {
		return Task{}, err
	}

	// Retrieve the updated task
	var task Task
	row := utils.DB.QueryRow("SELECT id, description,created_at,due_date,completed FROM tasks WHERE id = ?", id)
	err = row.Scan(&task.ID, &task.Description, &task.CreatedAt, &task.DueDate, &task.Completed)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}
