import React from "react";
import type { TreeNode, TreeNavigationProps } from "./types";
import styles from "./TreeNavigation.module.scss";

interface TreeNodeRendererProps {
  node: TreeNode;
  level: number;
  selectedNode: string | null;
  expandedNodes: Set<string>;
  onNodeSelect: (nodeId: string) => void;
  onNodeToggle: (nodeId: string, event: React.MouseEvent) => void;
}

const TreeNodeRenderer: React.FC<TreeNodeRendererProps> = ({
  node,
  level,
  selectedNode,
  expandedNodes,
  onNodeSelect,
  onNodeToggle,
}) => {
  const isExpanded = expandedNodes.has(node.id);
  const hasChildren = node.children && node.children.length > 0;

  return (
    <div
      key={node.id}
      className={styles.treeNode}
      style={{ "--level": `${level * 20}px` } as React.CSSProperties}
    >
      <div
        className={`${styles.treeNodeContent} ${
          selectedNode === node.id ? styles.selected : ""
        }`}
        onClick={() => onNodeSelect(node.id)}
      >
        {hasChildren && (
          <span
            className={styles.treeNodeExpander}
            onClick={(e) => onNodeToggle(node.id, e)}
          >
            {isExpanded ? "▼" : "▶"}
          </span>
        )}
        <span className={styles.treeNodeTitle}>{node.title}</span>
      </div>
      {isExpanded &&
        node.children?.map((child) => (
          <TreeNodeRenderer
            key={child.id}
            node={child}
            level={level + 1}
            selectedNode={selectedNode}
            expandedNodes={expandedNodes}
            onNodeSelect={onNodeSelect}
            onNodeToggle={onNodeToggle}
          />
        ))}
    </div>
  );
};

export const TreeNavigation: React.FC<TreeNavigationProps> = ({
  data,
  selectedNode,
  expandedNodes,
  onNodeSelect,
  onNodeToggle,
  className,
}) => {
  return (
    <div className={`${styles.treeNavigation} ${className || ""}`}>
      {data.map((node) => (
        <TreeNodeRenderer
          key={node.id}
          node={node}
          level={0}
          selectedNode={selectedNode}
          expandedNodes={expandedNodes}
          onNodeSelect={onNodeSelect}
          onNodeToggle={onNodeToggle}
        />
      ))}
    </div>
  );
};

export default TreeNavigation;
