import { useEffect, useState } from "react";
import { listTasks } from "@/client";
import type { TaskDigest, TaskStatus } from "@/client";
import { useNavigate } from "react-router-dom";
import { Flex, Card, Container, Progress } from "@gravity-ui/uikit";
import { Text } from "@gravity-ui/uikit";

export const TasksList = () => {
  const [tasks, setTasks] = useState<TaskDigest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
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

  if (loading) {
    return (
      <Container>
        <Text variant="header-1">Tasks</Text>
        <Text>Loading tasks...</Text>
      </Container>
    );
  }

  if (error) {
    return (
      <Container>
        <Text variant="header-1">Tasks</Text>
        <Text color="danger">{error}</Text>
      </Container>
    );
  }

  return (
    <Container>
      <Text variant="header-1">Tasks</Text>
      {tasks.length === 0 ? (
        <Text>No tasks found</Text>
      ) : (
        <Flex direction="column" gap="4">
          {tasks.map((task) => (
            <Card
              key={task.task_id}
              onClick={() => navigate(`/task/${task.task_id}`)}
              view="outlined"
              type='action'
              size="l"
            >
              <Container>
                <Flex justifyContent="space-between" alignItems="center">
                  <Text variant="header-2">#{task.task_id}</Text>
                  <Text variant="caption-1" color="secondary">
                    {task.triggered_by}
                  </Text>
                </Flex>
                <Text>{task.description}</Text>
                <Progress value={getProgressPercentage(task.status)} />
              </Container>
            </Card>
          ))}
        </Flex>
      )}
    </Container>
  );
};
