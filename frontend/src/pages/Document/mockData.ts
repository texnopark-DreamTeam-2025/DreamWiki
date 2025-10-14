/**
 * MOCK DATA - УДАЛИТЬ КОГДА БЭК ЗАРАБОТАЕТ
 *
 * Этот файл содержит все тестовые данные для разработки фронтенда
 * без подключения к бэкенду. Все константы и функции здесь - временные.
 *
 * TODO: Удалить весь файл после интеграции с API
 */

import type { V1DiagnosticInfoGetResponse } from "@/client";
import type { TreeNode } from "@/components/TreeNavigation";

// MOCK - Данные страницы для Document компонента
export const MOCK_PAGE_DATA: V1DiagnosticInfoGetResponse = {
  page: {
    page_id: "missouri-001",
    title: "Миссури",
    content: `
      <h2>История штата Миссури началась с эпохи заселения территории штата человеком в эпоху палеолитеза около 12 000 лет до нашей эры.</h2>

      <p>В конце XVII века эта земля была открыта французскими путешественниками <a href="#" style="color: blue; text-decoration: underline;">Жаком Маркеттом</a> и вошла в состав колонии <strong>Новая Франция</strong>. Берега <em>реки Миссури</em> заселялись мигрантами из <strong>Иллинойса</strong>, которые начали добычу <strong>свинца</strong> для экспорта в <strong>Европу</strong>.</p>

      <p>В 1803 году <strong>Франция</strong> продала <strong>Соединённые Штаты Америки</strong> все свои зап. территория <strong>Миссури</strong> стала частью приобретённой <strong>Территории Луизиана</strong>, переименованной в 1812 году в <strong>Территорию Миссури</strong>.</p>

      <p>В 1818 году территория подала запрос на вступление в США, что вызвало ожесточённые споры в Конгрессе и в итоге привело к принятию так называемого <strong>Миссурийского компромисса</strong>. В 1821 году <strong>Миссури</strong> был принят в <strong>США</strong> в качестве 24-го штата.</p>

      <h3>Основные периоды:</h3>
      <ul>
        <li><strong>Палеолит</strong> - первые поселения (12 000 лет до н.э.)</li>
        <li><strong>XVII век</strong> - открытие французскими путешественниками</li>
        <li><strong>1803 год</strong> - покупка Луизианы</li>
        <li><strong>1821 год</strong> - принятие в состав США</li>
      </ul>

      <blockquote>
        Штат быстро заселялся мигрантами из <strong>Германии</strong> и северных штатов, которые концентрировались в основном в городе <strong>Сент-Луис</strong>.
      </blockquote>

      <p>В годы <strong>Гражданской войны</strong> Миссури оставался под контролем <strong>Союза</strong>, однако на его территории происходило несколько сражений, а его жители служили как в армии Севера, так и в армии Юга.</p>

      <h3>Современный период</h3>
      <p>После <strong>Гражданской войны</strong> к власти в штате пришли республиканцы, но они расколились на радикальное и либеральное крыло, после чего к власти пришли начала <strong>либеральные республиканцы</strong>, а затем и демократы, которые в 1875 году приняли новую консервативную конституцию.</p>

      <p>Несмотря на строительство новых железных дорог, развитие промышленности и рост городов, <strong>Миссури</strong> по-прежнему отставал от ведущих индустриальных штатов. Когда началась <strong>Первая мировая война</strong>, штат получил лишь незначительную долю оборонных заказов.</p>
    `,
  },
};

// MOCK - Дерево навигации для TreeNavigation компонента
export const MOCK_TREE_DATA: TreeNode[] = [
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
          { id: "missouri-history", title: "История" },
          { id: "missouri-geography", title: "География" },
          { id: "missouri-economy", title: "Экономика" },
        ],
      },
      {
        id: "other-states",
        title: "Другие штаты",
        children: [
          { id: "california", title: "Калифорния" },
          { id: "texas", title: "Техас" },
          { id: "florida", title: "Флорида" },
        ],
      },
    ],
  },
  {
    id: "geography",
    title: "География",
    expanded: false,
    children: [
      {
        id: "usa",
        title: "США",
        children: [
          { id: "midwest", title: "Средний Запад" },
          { id: "south", title: "Юг" },
          { id: "west", title: "Запад" },
          { id: "northeast", title: "Северо-Восток" },
        ],
      },
      { id: "europe", title: "Европа" },
      { id: "asia", title: "Азия" },
    ],
  },
];

// MOCK - Начальные состояния для компонента Document
export const MOCK_INITIAL_SELECTED_NODE = "current";
export const MOCK_INITIAL_EXPANDED_NODES = new Set([
  "adyge-khabl",
  "logistics",
]);

// MOCK - Симуляция задержки API
export const MOCK_API_DELAY = 1000;

// MOCK - Функция симуляции загрузки данных
export const mockFetchPageData = async (
  id: string
): Promise<V1DiagnosticInfoGetResponse> => {
  // Симулируем задержку сети
  await new Promise((resolve) => setTimeout(resolve, MOCK_API_DELAY));

  console.log("🔄 MOCK: Загружаем данные для страницы ID:", id);

  return MOCK_PAGE_DATA;
};

// MOCK - Функция симуляции индексации страницы
export const mockIndexPage = async (pageId: string): Promise<void> => {
  console.log("🔄 MOCK: Индексируем страницу ID:", pageId);

  // Симулируем задержку обработки
  await new Promise((resolve) => setTimeout(resolve, 500));

  console.log("✅ MOCK: Страница успешно проиндексирована!");

  // В реальности здесь будет вызов indexatePage API
  // await indexatePage({ body: { page_id: pageId } });
};
