import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getTaskDetails } from "@/client";
import type { Task, TaskStatus } from "@/client";
import { Flex, Text, Button, Loader } from "@gravity-ui/uikit";

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
        return "positive";
      case "failed_by_error":
      case "failed_by_timeout":
        return "danger";
      case "cancelled":
        return "secondary";
      case "executing":
      default:
        return "info";
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
    return (
      <Flex justifyContent="center" alignItems="center" style={{ padding: 20 }}>
        <Loader size="m" />
        <Text variant="body-1" style={{ marginLeft: 10 }}>Loading task details...</Text>
      </Flex>
    );
  }

  if (error) {
    return (
      <Flex justifyContent="center" alignItems="center" style={{ padding: 20 }}>
        <Text variant="body-1" color="danger">{error}</Text>
      </Flex>
    );
  }

  if (!task) {
    return (
      <Flex justifyContent="center" alignItems="center" style={{ padding: 20 }}>
        <Text variant="body-1" color="danger">Task not found</Text>
      </Flex>
    );
  }

  return (
    <Flex direction="column" style={{ padding: 20 }}>
      <Flex justifyContent="space-between" alignItems="flex-start" style={{ marginBottom: 24, gap: 20 }}>
        <Flex direction="column">
          <Text variant="display-1">Task #{task.task_digest.task_id}</Text>
          <Flex
            alignItems="center"
            justifyContent="center"
            style={{
              backgroundColor: '#f0f0f0',
              padding: '4px 8px',
              borderRadius: 4,
              fontSize: 14,
              color: '#666666',
              display: 'inline-flex',
              marginTop: 8
            }}
          >
            {task.task_digest.triggered_by}
          </Flex>
        </Flex>
        <Button onClick={handleViewInternalState} size="m">
          View Internal State
        </Button>
      </Flex>

      <Flex direction="column" style={{ marginBottom: 24 }}>
        <Text variant="header-2" style={{ marginBottom: 12 }}>Description</Text>
        <Text variant="body-1" style={{ lineHeight: 1.5 }}>{task.task_digest.description}</Text>
      </Flex>

      <Flex direction="column">
        <Text variant="header-2" style={{ marginBottom: 16 }}>Subtasks</Text>
        {task.subtasks.length === 0 ? (
          <Text variant="body-1" color="secondary" style={{ fontStyle: 'italic' }}>No subtasks</Text>
        ) : (
          <Flex direction="column" style={{ gap: 16 }}>
            {task.subtasks.map((subtask, index) => (
              <Flex
                key={index}
                direction="column"
                style={{
                  border: '1px solid #e0e0e0',
                  borderRadius: 8,
                  padding: 16,
                  backgroundColor: '#ffffff',
                  boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)'
                }}
              >
                <Flex justifyContent="space-between" alignItems="center" style={{ marginBottom: 12 }}>
                  <Text variant="body-2">{subtask.description}</Text>
                  <Flex
                    alignItems="center"
                    justifyContent="center"
                    style={{
                      padding: '4px 8px',
                      borderRadius: 4,
                      fontSize: 12,
                      fontWeight: 500
                    }}
                  >
                    <Text variant="caption-1" color={getSubtaskStatusColor(subtask.status)}>
                      {getSubtaskStatusText(subtask.status)}
                    </Text>
                  </Flex>
                </Flex>
                {subtask.subsubtasks.length > 0 && (
                  <Flex
                    direction="column"
                    style={{
                      marginTop: 12,
                      paddingLeft: 16,
                      borderLeft: '2px solid #f0f0f0'
                    }}
                  >
                    {subtask.subsubtasks.map((subsubtask, subIndex) => (
                      <Flex
                        key={subIndex}
                        justifyContent="space-between"
                        alignItems="center"
                        style={{
                          padding: '8px 0',
                          borderBottom: '1px solid #f0f0f0'
                        }}
                      >
                        <Text variant="body-1" color="secondary" style={{ fontSize: 14 }}>
                          {subsubtask.description}
                        </Text>
                        <Flex
                          alignItems="center"
                          justifyContent="center"
                          style={{
                            padding: '2px 6px',
                            borderRadius: 4,
                            fontSize: 11,
                            fontWeight: 500
                          }}
                        >
                          <Text variant="caption-2" color={getSubtaskStatusColor(subsubtask.status)}>
                            {getSubtaskStatusText(subsubtask.status)}
                          </Text>
                        </Flex>
                      </Flex>
                    ))}
                  </Flex>
                )}
              </Flex>
            ))}
          </Flex>
        )}
      </Flex>
    </Flex>
  );
};
