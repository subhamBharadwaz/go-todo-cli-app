/*
Copyright © 2024 Subham Bharadwaz subhamsbharadwaz@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
	"github.com/subhamBharadwaz/go-todo-cli-app/utils"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new todo item to the CSV file",
	Long: `The add command allows you to add a new todo item to the CSV file 'tasks.csv'. You need to provide a description for the todo item. The command will generate a unique ID for the new task and record the current date and time as the creation time.

Example usage:
  $ tasks add "Buy groceries"

This command updates the 'tasks.csv' file with the new todo item, including its ID, description, and creation time. The task is appended to the end of the file, and the newly added task is displayed in a tabular format upon successful addition.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("Please provide a task description.")
			return
		}

		description := args[0]
		dueDate, _ := cmd.Flags().GetString("due")

		task, err := addTask(description, dueDate)
		if err != nil {
			fmt.Println("Error adding task:", err)
		} else {

			writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
			fmt.Fprintln(writer, "ID\tName\tCreated\tDue\t")

			createdTime, _ := time.Parse(time.RFC3339, task.CreatedAt)

			fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t\n", task.ID, task.Description, timediff.TimeDiff(createdTime), dueDate)

			writer.Flush()
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("due", "", "Specify the due date for the task")

}

type Task struct {
	ID          int
	Description string
	CreatedAt   string
	DueDate     string
	Completed   bool
}

func addTask(description string, dueDate string) (Task, error) {
	if utils.DB == nil {
		return Task{}, fmt.Errorf("database is not initialized")
	}

	// Get the current time
	createdAt := time.Now().Format(time.RFC3339)

	result, err := utils.DB.Exec("INSERT INTO tasks (description, created_at, due_date, completed) VALUES (?, ?, ?, ?)", description, createdAt, dueDate, false)
	if err != nil {
		return Task{}, fmt.Errorf("failed to insert task: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	task := Task{
		ID:          int(id),
		Description: description,
		CreatedAt:   createdAt,
		DueDate:     dueDate,
	}

	return task, nil
}
