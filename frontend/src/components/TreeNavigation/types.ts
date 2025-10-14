export interface TreeNode {
  id: string;
  title: string;
  children?: TreeNode[];
  expanded?: boolean;
}

export interface TreeNavigationProps {
  /** Данные дерева */
  data: TreeNode[];
  /** ID выбранного узла */
  selectedNode: string | null;
  /** Множество раскрытых узлов */
  expandedNodes: Set<string>;
  /** Коллбек для выбора узла */
  onNodeSelect: (nodeId: string) => void;
  /** Коллбек для переключения раскрытия узла */
  onNodeToggle: (nodeId: string, event: React.MouseEvent) => void;
  /** Дополнительный CSS класс */
  className?: string;
}
