import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SearchDone from "@/pages/SearchDone";
import Search from "./modules/Search/Search";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<div>Главная страница</div>} />
        <Route path="/search-done" element={<Search />} />
        {/* Добавьте другие маршруты здесь */}
        <Route path="*" element={<div>Страница не найдена</div>} />
      </Routes>
    </Router>
  );
}

export default App;
