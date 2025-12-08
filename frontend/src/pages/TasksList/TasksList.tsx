import { useInfiniteQuery } from '@tanstack/react-query';
import { listTasks } from "@/client";
import type { TaskDigest, TaskStatus } from "@/client";
import { useNavigate } from "react-router-dom";
import { Flex, Card, Progress, Text, Loader, Label, Box } from "@gravity-ui/uikit";
import { useCallback, useRef, useEffect } from 'react';

const TASKS_PER_PAGE = 20;

export const TasksList = () => {
  const navigate = useNavigate();
  const observer = useRef<IntersectionObserver | null>(null);

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    isError
  } = useInfiniteQuery({
    queryKey: ['tasks'],
    queryFn: async ({ pageParam }) => {
      const response = await listTasks({
        body: {
          filters: {
            only_my_tasks: false,
            only_active: false,
          },
          cursor: pageParam
        },
      });

      return response.data;
    },
    getNextPageParam: (lastPage) => {
      if (lastPage?.next_info?.has_more) {
        return lastPage.next_info.cursor;
      }
      return undefined;
    },
    initialPageParam: undefined as string | undefined,
  });

  const lastTaskElementRef = useCallback((node: HTMLDivElement | null) => {
    if (isLoading || isFetchingNextPage) return;
    if (observer.current) observer.current.disconnect();

    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasNextPage) {
        fetchNextPage();
      }
    });

    if (node) observer.current.observe(node);
  }, [isLoading, isFetchingNextPage, hasNextPage, fetchNextPage]);

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

  const tasks = data?.pages.flatMap(page => page?.tasks || []) || [];

  return (
    <Flex direction="column" gap="4">
      <Text variant="header-1">Задачи</Text>
      {isLoading && <Loader />}
      {isError && <Text>Error loading tasks</Text>}
      {!isLoading && tasks.length === 0 && <Text>No tasks found</Text>}
      {!isLoading && tasks.length > 0 && (
        <Flex direction="column" gap="4">
          {tasks.map((task, index) => {
            if (tasks.length === index + 1) {
              return (
                <div ref={lastTaskElementRef} key={task.task_id}>
                  <Card
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
                </div>
              );
            } else {
              return (
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
              );
            }
          })}
          {isFetchingNextPage && (
            <Flex justifyContent="center" alignItems="center" style={{ padding: '20px' }}>
              <Loader />
            </Flex>
          )}
        </Flex>
      )}
    </Flex>
  );
};
