import { useEffect, useState } from "react";
import { useNavigate, useParams, Link } from "react-router-dom";
import { getTaskDetails } from "@/client";
import type { Task, TaskStatus } from "@/client";
import { Flex, Text, Button, Loader, Card, Label, Breadcrumbs, Box } from "@gravity-ui/uikit";
import { ActionBar } from "@gravity-ui/navigation";

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
      <Flex justifyContent="center" alignItems="center" gap="2">
        <Loader size="m" />
        <Text variant="body-1">Loading task details...</Text>
      </Flex>
    );
  }

  if (error) {
    return (
      <Flex justifyContent="center" alignItems="center">
        <Text variant="body-1" color="danger">
          {error}
        </Text>
      </Flex>
    );
  }

  if (!task) {
    return (
      <Flex justifyContent="center" alignItems="center">
        <Text variant="body-1" color="danger">
          Task not found
        </Text>
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="5">
      <ActionBar>
        <Flex alignItems="center" dir="horizontal" width="100%">
          <Breadcrumbs style={{ width: "100%" }} showRoot>
            <Breadcrumbs.Item onClick={() => navigate("/tasks")}>Задачи</Breadcrumbs.Item>
            <Breadcrumbs.Item>Задача #{task.task_digest.task_id}</Breadcrumbs.Item>
          </Breadcrumbs>
        </Flex>
      </ActionBar>

      <Flex justifyContent="space-between" alignItems="flex-start" gap="4">
        <Flex direction="column" gap="2">
          <Text variant="display-1">Task #{task.task_digest.task_id}</Text>
          <Label theme="normal" size="m">
            {task.task_digest.triggered_by}
          </Label>
        </Flex>
        <Button onClick={handleViewInternalState} size="m">
          View Internal State
        </Button>
      </Flex>

      <Card theme="normal" size="l">
        <Flex direction="column" gap="3">
          <Text variant="header-2">Description</Text>
          <Text variant="body-1">{task.task_digest.description}</Text>
        </Flex>
      </Card>

      <Flex direction="column" gap="4">
        <Text variant="header-2">Subtasks</Text>
        {task.subtasks.length === 0 ? (
          <Text variant="body-1" color="secondary">
            No subtasks
          </Text>
        ) : (
          <Flex direction="column" gap="4">
            {task.subtasks.map((subtask, index) => (
              <Card key={index} theme="normal" size="l">
                <Flex direction="column" gap="3">
                  <Flex justifyContent="space-between" alignItems="center">
                    <Text variant="body-2">{subtask.description}</Text>
                    <Label theme={getSubtaskStatusColor(subtask.status) as any} size="s">
                      {getSubtaskStatusText(subtask.status)}
                    </Label>
                  </Flex>
                  {subtask.subsubtasks.length > 0 && (
                    <Flex direction="column" gap="2" style={{ paddingLeft: 16 }}>
                      {subtask.subsubtasks.map((subsubtask, subIndex) => (
                        <Flex
                          key={subIndex}
                          justifyContent="space-between"
                          alignItems="center"
                          gap="4"
                        >
                          <Text variant="body-1" color="secondary">
                            {subsubtask.description}
                          </Text>
                          <Label theme={getSubtaskStatusColor(subsubtask.status) as any} size="xs">
                            {getSubtaskStatusText(subsubtask.status)}
                          </Label>
                        </Flex>
                      ))}
                    </Flex>
                  )}
                </Flex>
              </Card>
            ))}
          </Flex>
        )}
      </Flex>
    </Flex>
  );
};
