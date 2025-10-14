import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { AsideHeader, type MenuItem, FooterItem } from "@gravity-ui/navigation";
import {
  House,
  Magnifier,
  Gear,
  Clock,
  FileText,
  Square,
  Person,
  Plus,
} from "@gravity-ui/icons";
import styles from "./AsideBar.module.scss";

interface AsideBarProps {
  children: React.ReactNode;
}

export default function AsideBar({ children }: AsideBarProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const [compact, setCompact] = useState(false);

  const logo = {
    text: "DreamWiki",
    icon: House,
    onClick: () => navigate("/"),
  };

  const subheaderItems = [
    {
      item: {
        id: "services",
        title: "Все сервисы",
        icon: Square,
        current: location.pathname === "/services",
        onItemClick: () => navigate("/services"),
        qa: "services",
        tooltipText: "Просмотр всех доступных сервисов",
      },
      enableTooltip: true,
    },
    {
      item: {
        id: "search",
        title: "Поиск",
        icon: Magnifier,
        current:
          location.pathname === "/search" ||
          location.pathname === "/search-done",
        onItemClick: () => navigate("/search"),
        qa: "search",
        tooltipText: "Поиск по документам и статьям",
      },
      enableTooltip: true,
    },
  ];

  const menuItems: MenuItem[] = [
    {
      id: "integrations",
      title: "Настройка интеграций",
      icon: Gear,
      current: location.pathname === "/integrations",
      onItemClick: () => navigate("/integrations"),
      category: "Настройки",
      qa: "integrations",
      tooltipText: "Управление интеграциями с внешними сервисами",
    },
    {
      id: "logs",
      title: "Журнал интеграций",
      icon: Clock,
      current: location.pathname === "/logs",
      onItemClick: () => navigate("/logs"),
      category: "Настройки",
      qa: "logs",
      tooltipText: "История операций интеграций",
    },
    {
      id: "drafts",
      title: "Черновики",
      icon: FileText,
      current: location.pathname === "/drafts",
      onItemClick: () => navigate("/drafts"),
      category: "Контент",
      qa: "drafts",
      tooltipText: "Управление черновиками документов",
    },
    {
      id: "divider",
      title: "",
      type: "divider",
    },
    {
      id: "create-draft",
      title: "Создать черновик",
      icon: Plus,
      type: "action",
      onItemClick: () => navigate("/drafts/new"),
      category: "Действия",
      qa: "create-draft-button",
      tooltipText: "Создать новый черновик документа",
    },
  ];

  const renderFooter = () => (
    <>
      <FooterItem
        item={{
          id: "settings",
          title: "Настройки",
          icon: Gear,
          onItemClick: () => navigate("/settings"),
        }}
        compact={compact}
        bringForward={true}
      />
      <FooterItem
        item={{
          id: "account",
          title: "Account",
          icon: Person,
          onItemClick: () => navigate("/profile"),
        }}
        compact={compact}
        bringForward={false}
      />
    </>
  );

  const renderContent = () => (
    <main className={styles.asideBarMain}>{children}</main>
  );

  return (
    <div className={`${styles.asideBarContainer} ${styles.asideBar}`}>
      <AsideHeader
        compact={compact}
        onChangeCompact={setCompact}
        logo={logo}
        subheaderItems={subheaderItems}
        menuItems={menuItems}
        headerDecoration
        renderFooter={renderFooter}
        renderContent={renderContent}
        collapseTitle="Свернуть"
        expandTitle="Развернуть"
        menuMoreTitle="Ещё"
        multipleTooltip={true}
        qa="dreamwiki-aside-header"
        customBackground={
          <div
            style={{
              background: "linear-gradient(180deg, #1f2937, #111827)",
              height: "100%",
            }}
          />
        }
        customBackgroundClassName={styles.customBackground}
      />
    </div>
  );
}
