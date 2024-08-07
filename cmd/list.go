/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"strconv"

	"time"

	"github.com/subhamBharadwaz/go-todo-cli-app/utils"

	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all todo items in the CSV file",
	Long: `The list command displays all todo items currently stored in the CSV file 'tasks.csv'. Each todo item is shown with its ID, description, and the time it was created. This command is useful for reviewing all the tasks you have added.

Example usage:
  $ tasks list

This command reads from the 'tasks.csv' file and formats the output in a tabular format for easy reading. If the file is empty or contains no todos, it will indicate that there are no tasks available. But this command will gives only the todos which done status is false.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the -all flag is set
		showAll, _ := cmd.Flags().GetBool("all")

		tasks, err := getAllTasks()
		if err != nil {
			fmt.Println("Error fetching tasks:", err)
			return
		}

		var columns []string

		// Print headers based on --all flag
		if showAll {
			columns = append(columns, "ID", "Name", "Created", "Due", "Done")
		} else {
			columns = append(columns, "ID", "Name", "Created", "Due")
		}

		var rows [][]string

		for _, task := range tasks {
			// Print task details based on --all flag
			createdTime, _ := time.Parse(time.RFC3339, task.CreatedAt)
			completed := strconv.FormatBool(task.Completed)

			if showAll {
				// fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\t\n", task.ID, task.Description, timediff.TimeDiff(createdTime), task.DueDate, completed)
				rows = append(rows, []string{
					fmt.Sprintf("%d", task.ID),
					task.Description,
					timediff.TimeDiff(createdTime),
					task.DueDate,
					completed,
				})

			} else if completed == "false" {

				rows = append(rows, []string{
					fmt.Sprintf("%d", task.ID),
					task.Description,
					timediff.TimeDiff(createdTime),
					task.DueDate,
				})

			}

		}

		t := utils.CreateStyledTable(columns, rows)
		fmt.Println(t)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "Show all tasks, including completed ones")
}

func getAllTasks() ([]Task, error) {
	rows, err := utils.DB.Query("SELECT id, description, created_at, due_date, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Description, &task.CreatedAt, &task.DueDate, &task.Completed)

		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)

	}
	return tasks, nil
}
