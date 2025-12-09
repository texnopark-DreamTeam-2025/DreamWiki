import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getTaskInternalState } from "@/client";
import type { RawJson } from "@/client";
import { MonacoEditor } from "@/components/MonacoEditor/MonacoEditor";
import { ActionBar } from "@gravity-ui/navigation";
import { Breadcrumbs, Flex, Text } from "@gravity-ui/uikit";

export const InternalState = () => {
  const { taskId } = useParams<{ taskId: string }>();
  const navigate = useNavigate();
  const [taskState, setTaskState] = useState<RawJson | null>(null);
  const [actions, setActions] = useState<RawJson[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchInternalState = async () => {
      if (!taskId) {
        setError("Task ID is missing");
        setLoading(false);
        return;
      }

      try {
        setLoading(true);
        const response = await getTaskInternalState({
          body: {
            task_id: parseInt(taskId, 10),
          },
        });

        setTaskState(response.data?.task_state || null);
        setActions(response.data?.actions || []);
      } catch (err) {
        setError("Failed to fetch internal state");
        console.error("Error fetching internal state:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchInternalState();
  }, [taskId]);

  if (loading) {
    return (
      <Flex direction="column" alignItems="center" justifyContent="center" style={{ height: "100vh" }}>
        <Text variant="body-2">Loading internal state...</Text>
      </Flex>
    );
  }

  if (error) {
    return (
      <Flex direction="column" alignItems="center" justifyContent="center" style={{ height: "100vh" }}>
        <Text variant="body-2" color="danger">
          {error}
        </Text>
      </Flex>
    );
  }

  return (
    <Flex direction="column" gap="4">
      <ActionBar>
        <Flex alignItems="center" direction="row" width="100%">
          <Breadcrumbs style={{ width: "100%" }} showRoot>
            <Breadcrumbs.Item href="/tasks" onClick={() => navigate("/tasks")}>Задачи</Breadcrumbs.Item>
            <Breadcrumbs.Item href={`/task/${taskId}`} onClick={() => navigate(`/task/${taskId}`)}>Задача #{taskId}</Breadcrumbs.Item>
            <Breadcrumbs.Item>Internal State</Breadcrumbs.Item>
          </Breadcrumbs>
        </Flex>
      </ActionBar>

      <Text variant="header-1">Internal State for Task #{taskId}</Text>

      <MonacoEditor
        value={JSON.stringify({ task_state: taskState, actions }, null, 2)}
        language="json"
        height="600px"
        readOnly={true}
      />
    </Flex>
  );
};
