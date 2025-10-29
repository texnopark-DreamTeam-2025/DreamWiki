import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getTaskDetails } from "@/client";
import type { Task, TaskStatus } from "@/client";
import styles from "./TaskDetails.module.scss";

export const TaskDetails = () => {
  const { taskId } = useParams<{ taskId: string }>();
  const navigate = useNavigate();
  const [task, setTask] = useState<Task | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTask = async () => {
      if (!taskId) {
        setError("Task ID is missing");
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        const response = await getTaskDetails({
          body: {
            task_id: parseInt(taskId, 10),
          },
        });
        setTask(response.data?.task || null);
      } catch (err) {
        setError("Failed to fetch task details");
        console.error("Error fetching task details:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchTask();
  }, [taskId]);

  const handleViewInternalState = async () => {
    if (!taskId) return;

    try {
      // Navigate to internal state page
      navigate(`/task/${taskId}/internal-state`);
    } catch (err) {
      console.error("Error navigating to internal state:", err);
    }
  };

  const getSubtaskStatusColor = (status: TaskStatus) => {
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

  const getSubtaskStatusText = (status: TaskStatus) => {
    switch (status) {
      case "done":
        return "Done";
      case "failed_by_error":
        return "Failed (Error)";
      case "failed_by_timeout":
        return "Failed (Timeout)";
      case "cancelled":
        return "Cancelled";
      case "executing":
      default:
        return "Executing";
    }
  };

  if (loading) {
    return <div className={styles.loading}>Loading task details...</div>;
  }

  if (error) {
    return <div className={styles.error}>{error}</div>;
  }

  if (!task) {
    return <div className={styles.error}>Task not found</div>;
  }

  return (
    <div className={styles.taskDetails}>
      <div className={styles.header}>
        <div className={styles.taskInfo}>
          <h1>Task #{task.task_digest.task_id}</h1>
          <div className={styles.taskTrigger}>{task.task_digest.triggered_by}</div>
        </div>
        <button className={styles.internalStateButton} onClick={handleViewInternalState}>
          View Internal State
        </button>
      </div>

      <div className={styles.description}>
        <h2>Description</h2>
        <p>{task.task_digest.description}</p>
      </div>

      <div className={styles.subtasks}>
        <h2>Subtasks</h2>
        {task.subtasks.length === 0 ? (
          <div className={styles.noSubtasks}>No subtasks</div>
        ) : (
          <div className={styles.subtasksList}>
            {task.subtasks.map((subtask, index) => (
              <div key={index} className={styles.subtaskItem}>
                <div className={styles.subtaskHeader}>
                  <div className={styles.subtaskDescription}>{subtask.description}</div>
                  <div
                    className={`${styles.subtaskStatus} ${getSubtaskStatusColor(subtask.status)}`}
                  >
                    {getSubtaskStatusText(subtask.status)}
                  </div>
                </div>
                {subtask.subsubtasks.length > 0 && (
                  <div className={styles.subsubtasks}>
                    {subtask.subsubtasks.map((subsubtask, subIndex) => (
                      <div key={subIndex} className={styles.subsubtaskItem}>
                        <div className={styles.subsubtaskDescription}>{subsubtask.description}</div>
                        <div
                          className={`${styles.subsubtaskStatus} ${getSubtaskStatusColor(subsubtask.status)}`}
                        >
                          {getSubtaskStatusText(subsubtask.status)}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
