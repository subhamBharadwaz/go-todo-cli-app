/*
Copyright © 2024 Subham Bharadwaz subhamsbharadwaz@gmail.com
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
	"todo-list/utils"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
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

		if len(args) == 0 {
			fmt.Println("Please provide a task description.")
			return
		}

		description := args[0]
		task, err := addTask(description)
		if err != nil {
			fmt.Println("Error adding task:", err)
		} else {

			writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
			fmt.Fprintln(writer, "ID\tName\tCreated\t")

			createdTime, _ := time.Parse(time.RFC3339, task[2])

			fmt.Fprintf(writer, "%s\t%s\t%s\t\n", task[0], task[1], timediff.TimeDiff(createdTime))

			writer.Flush()
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

}

func addTask(description string) ([]string, error) {
	file, err := utils.LoadFile("tasks.csv")
	if err != nil {
		return nil, err
	}
	defer utils.CloseFile(file)

	// Read existing tasks to determine the next ID
	existingTasks, err := utils.ReadTasks(file)
	if err != nil {
		return nil, err
	}

	// Determine the next ID
	nextID := 1
	if len(existingTasks) > 0 {
		lastTask := existingTasks[len(existingTasks)-1]
		lastID, err := strconv.Atoi(lastTask[0])
		if err != nil {
			return nil, err
		}
		nextID = lastID + 1
	}

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Get the current time
	createdTime := time.Now().Format(time.RFC3339)

	// task done
	complete := "false"

	// Write the task description to the CSV file
	task := []string{strconv.Itoa(nextID), description, createdTime, complete}
	if err := writer.Write(task); err != nil {
		return nil, err
	}

	return task, nil
}
