import { useEffect, useState } from "react";
import { listTasks } from "@/client";
import type { TaskDigest, TaskStatus } from "@/client";
import styles from "./TasksList.module.scss";

export const TasksList = () => {
  const [tasks, setTasks] = useState<TaskDigest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        setLoading(true);
        const response = await listTasks({
          body: {
            filters: {
              only_my_tasks: false,
              only_active: false,
            },
          },
        });
        setTasks(response.data?.tasks || []);
      } catch (err) {
        setError("Failed to fetch tasks");
        console.error("Error fetching tasks:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, []);

  const getProgressPercentage = (status: TaskStatus) => {
    switch (status) {
      case "done":
        return 100;
      case "failed_by_error":
      case "failed_by_timeout":
      case "cancelled":
        return 0;
      case "executing":
      default:
        return 50;
    }
  };

  const getStatusColor = (status: TaskStatus) => {
    switch (status) {
      case "done":
        return styles.statusDone;
      case "failed_by_error":
      case "failed_by_timeout":
        return styles.statusFailed;
      case "cancelled":
        return styles.statusCancelled;
      case "executing":
      default:
        return styles.statusExecuting;
    }
  };

  if (loading) {
    return <div className={styles.loading}>Loading tasks...</div>;
  }

  if (error) {
    return <div className={styles.error}>{error}</div>;
  }

  return (
    <div className={styles.tasksList}>
      <h1>Tasks</h1>
      {tasks.length === 0 ? (
        <div className={styles.noTasks}>No tasks found</div>
      ) : (
        <div className={styles.taskList}>
          {tasks.map((task) => (
            <div key={task.task_id} className={styles.taskItem}>
              <div className={styles.taskHeader}>
                <div className={styles.taskNumber}>#{task.task_id}</div>
                <div className={styles.taskTrigger}>{task.triggered_by}</div>
              </div>
              <div className={styles.taskDescription}>{task.description}</div>
              <div className={styles.progressBarContainer}>
                <div className={styles.progressBar}>
                  <div
                    className={`${styles.progressFill} ${getStatusColor(task.status)}`}
                    style={{ width: `${getProgressPercentage(task.status)}%` }}
                  ></div>
                </div>
                <div className={styles.progressText}>{getProgressPercentage(task.status)}%</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
