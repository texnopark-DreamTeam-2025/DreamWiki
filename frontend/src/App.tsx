import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SearchDone from "@/pages/SearchDone";
import Search from "./modules/Search/Search";
import PageSearch from "./pages/PageSearch/PageSearch";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<div>Главная страница</div>} />
        <Route path="/search-done" element={<SearchDone/>} />
        {/* Добавьте другие маршруты здесь */}
        <Route path="*" element={<div>Страница не найдена</div>} />
      </Routes>
    </Router>
  );
}

export default App;
