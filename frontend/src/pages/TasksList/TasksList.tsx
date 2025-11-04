import { useEffect, useState } from "react";
import { listTasks } from "@/client";
import type { TaskDigest, TaskStatus } from "@/client";
import { useNavigate } from "react-router-dom";
import { Flex, Card, Progress, Text, Loader, Label, Box } from "@gravity-ui/uikit";

export const TasksList = () => {
  const [tasks, setTasks] = useState<TaskDigest[] | undefined>();
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

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
        console.error("Error fetching tasks:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, []);

  const getTaskStatusColor = (status: TaskStatus) => {
    switch (status) {
      case "done":
        return "success";
      case "failed_by_error":
      case "failed_by_timeout":
        return "danger";
      case "cancelled":
        return "misc";
      case "executing":
      default:
        return "info";
    }
  };

  const getTaskStatusText = (status: TaskStatus) => {
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

  return (
    <Flex direction="column" gap="4">
      <Text variant="header-1">Задачи</Text>
      {loading && <Loader />}
      {!loading && !tasks && <Text>Error loading tasks</Text>}
      {!loading && tasks && tasks.length === 0 && <Text>No tasks found</Text>}
      {!loading && tasks && tasks.length > 0 && (
        <Flex direction="column" gap="4">
          {tasks.map((task) => (
            <Card
              key={task.task_id}
              theme="normal"
              size="l"
              style={{ cursor: "pointer", padding: 10 }}
            >
              <Flex direction="column" gap="3" onClick={() => navigate(`/task/${task.task_id}`)}>
                <Flex justifyContent="space-between" alignItems="center">
                  <Text variant="header-2">#{task.task_id}</Text>
                  <Label theme={getTaskStatusColor(task.status) as any} size="m">
                    {getTaskStatusText(task.status)}
                  </Label>
                </Flex>
                <Text variant="body-1">{task.description}</Text>
                <Box width="100%">
                  <Progress
                    colorStops={[
                      {
                        stop: task.progress_percentage,
                        theme: getTaskStatusColor(task.status),
                      },
                    ]}
                    value={task.progress_percentage}
                  />
                </Box>
                <Text variant="body-2" color="secondary">
                  Triggered by: {task.triggered_by}
                </Text>
              </Flex>
            </Card>
          ))}
        </Flex>
      )}
    </Flex>
  );
};
