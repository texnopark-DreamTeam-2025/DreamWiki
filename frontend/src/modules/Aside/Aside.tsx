import { useState } from "react";
import { AsideHeader, type MenuItem, type Logo } from "@gravity-ui/navigation";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useNavigate,
} from "react-router-dom";
import { Home, Search, Settings } from "@gravity-ui/icons"; // пример иконок

function AppContent() {
  const navigate = useNavigate();
  const [compact, setCompact] = useState(false);

  // 👇 Пример логотипа
  const logo: Logo = {
    text: "MyApp",
    icon: <Home />,
    onClick: () => navigate("/"),
  };

  // 👇 Пример пунктов меню
  const menuItems: MenuItem[] = [
    {
      id: "home",
      title: "Главная",
      icon: <Home />,
      current: window.location.pathname === "/",
      onItemClick: () => navigate("/"),
    },
    {
      id: "search",
      title: "Поиск",
      icon: <Search />,
      current: window.location.pathname === "/search",
      onItemClick: () => navigate("/search"),
    },
    {
      id: "settings",
      title: "Настройки",
      icon: <Settings />,
      current: window.location.pathname === "/settings",
      onItemClick: () => navigate("/settings"),
    },
  ];

  // 👇 Нижний блок футера
  const renderFooter = () => (
    <div style={{ padding: "10px", fontSize: "13px", color: "#999" }}>
      © 2025 MyApp
    </div>
  );

  return (
    <AsideHeader
      compact={compact}
      onChangeCompact={setCompact}
      logo={logo}
      menuItems={menuItems}
      headerDecoration
      renderFooter={renderFooter}
      customBackground={
        <div
          style={{
            background: "linear-gradient(180deg, #1f2937, #111827)",
            height: "100%",
          }}
        />
      }
      customBackgroundClassName="aside-background"
      collapseTitle="Свернуть"
      expandTitle="Развернуть"
    >
      {/* Контент страницы */}
      <div style={{ padding: "20px" }}>
        <Routes>
          <Route path="/" element={<h2>Главная страница</h2>} />
          <Route path="/search" element={<h2>Страница поиска</h2>} />
          <Route path="/settings" element={<h2>Настройки</h2>} />
        </Routes>
      </div>
    </AsideHeader>
  );
}

export default function AsideHeaderExample() {
  return (
    <Router>
      <AppContent />
    </Router>
  );
}
