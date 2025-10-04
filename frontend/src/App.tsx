import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SearchDone from "@/pages/SearchDone";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<div>Главная страница</div>} />
        <Route path="/search-done" element={<SearchDone />} />
        {/* Добавьте другие маршруты здесь */}
        <Route path="*" element={<div>Страница не найдена</div>} />
      </Routes>
    </Router>
  );
}

export default App;
