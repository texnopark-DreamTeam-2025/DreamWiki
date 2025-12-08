import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getTaskInternalState } from "@/client";
import type { RawJson } from "@/client";
import { MonacoEditor } from "@/components/MonacoEditor/MonacoEditor";
import styles from "./InternalState.module.scss";

export const InternalState = () => {
  const { taskId } = useParams<{ taskId: string }>();
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
    return <div className={styles.loading}>Loading internal state...</div>;
  }

  if (error) {
    return <div className={styles.error}>{error}</div>;
  }

  return (
    <div className={styles.internalState}>
      <h1>Internal State for Task #{taskId}</h1>

      <div className={styles.section}>
        <h2>Task State and Actions</h2>
        <MonacoEditor
          value={JSON.stringify({ task_state: taskState, actions }, null, 2)}
          language="json"
          height="600px"
          readOnly={true}
        />
      </div>
    </div>
  );
};
