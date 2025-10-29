import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import {
  AsideHeader,
  type MenuItem,
  FooterItem,
  type SubheaderMenuItem,
} from "@gravity-ui/navigation";
import {
  House,
  Magnifier,
  Gear,
  Clock,
  FileText,
  Person,
  Plus,
  ArrowRightFromSquare,
} from "@gravity-ui/icons";
import styles from "./AsideBar.module.scss";
import { useAuth } from "@/contexts";

interface AsideBarProps {
  children: React.ReactNode;
}

export default function AsideBar({ children }: AsideBarProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const [compact, setCompact] = useState(false);
  const { isAuthenticated, logout } = useAuth();

  const logo = {
    text: "DreamWiki",
    icon: House,
    onClick: () => navigate("/"),
  };

  const subheaderItems: SubheaderMenuItem[] = [
    {
      item: {
        id: "search",
        title: "Поиск",
        icon: Magnifier,
        current: location.pathname === "/search" || location.pathname === "/search-done",
        onItemClick: () => navigate("/search"),
        qa: "search",
        tooltipText: "Поиск по документам и статьям",
      },
      enableTooltip: true,
    },
  ];

  const menuItems: MenuItem[] = isAuthenticated
    ? [
        {
          id: "integrations",
          title: "Настройка интеграций",
          icon: Gear,
          current: location.pathname === "/integration-settings",
          onItemClick: () => navigate("/integration-settings"),
          category: "Настройки",
          qa: "integrations",
          tooltipText: "Управление интеграциями с внешними сервисами",
        },
        {
          id: "logs",
          title: "Журнал интеграций",
          icon: Clock,
          current: location.pathname === "/integration-logs",
          onItemClick: () => navigate("/integration-logs"),
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
          id: "tasks",
          title: "Задачи",
          icon: Clock,
          current: location.pathname === "/tasks",
          onItemClick: () => navigate("/tasks"),
          category: "Система",
          qa: "tasks",
          tooltipText: "Просмотр и управление задачами",
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
      ]
    : [
        {
          id: "authorization",
          title: "Авторизация",
          icon: Person,
          current: location.pathname === "/login",
          onItemClick: () => navigate("/login"),
          category: "Аккаунт",
          qa: "authorization",
          tooltipText: "Вход в систему",
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
      {isAuthenticated ? (
        <FooterItem
          item={{
            id: "logout",
            title: "Выйти",
            icon: ArrowRightFromSquare,
            onItemClick: () => {
              logout();
              navigate("/login");
            },
          }}
          compact={compact}
          bringForward={false}
        />
      ) : (
        <FooterItem
          item={{
            id: "account",
            title: "Account",
            icon: Person,
            onItemClick: () => navigate("/login"),
          }}
          compact={compact}
          bringForward={false}
        />
      )}
    </>
  );

  const renderContent = () => <main className={styles.asideBarMain}>{children}</main>;

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
