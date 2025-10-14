/**
 * Данные для дерева навигации
 *
 * TODO: Заменить на реальные данные из API
 */

import type { TreeNode } from "@/components/TreeNavigation";

// Дерево навигации для TreeNavigation компонента
export const TREE_DATA: TreeNode[] = [
  {
    id: "logistics",
    title: "Логистика",
    expanded: true,
    children: [
      {
        id: "warehouses-sub",
        title: "Склады",
        children: [
          { id: "warehouse-moscow", title: "Склад Москва" },
          { id: "warehouse-spb", title: "Склад СПб" },
        ],
      },
      {
        id: "transport",
        title: "Транспорт",
        children: [
          { id: "trucks", title: "Грузовики" },
          { id: "trains", title: "Железная дорога" },
        ],
      },
    ],
  },
  {
    id: "warehouses",
    title: "Склады",
    expanded: false,
    children: [
      { id: "moscow", title: "Москва" },
      { id: "spb", title: "Санкт-Петербург" },
      { id: "kazan", title: "Казань" },
      { id: "novosibirsk", title: "Новосибирск" },
    ],
  },
  {
    id: "adyge-khabl",
    title: "Адыге-Хабль",
    expanded: true,
    children: [
      {
        id: "current",
        title: "Миссури",
        children: [
          { id: "mimi", title: "Мими" },
          { id: "history", title: "История" },
          { id: "geography", title: "География" },
        ],
      },
      {
        id: "other",
        title: "Другие документы",
        children: [
          { id: "doc1", title: "Документ 1" },
          { id: "doc2", title: "Документ 2" },
        ],
      },
    ],
  },
];

// Изначально выбранный узел
export const INITIAL_SELECTED_NODE = "mimi";

// Изначально раскрытые узлы
export const INITIAL_EXPANDED_NODES = new Set(["adyge-khabl", "logistics"]);
