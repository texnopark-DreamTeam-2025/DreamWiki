import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AsideBar from "@/components/AsideBar";
import SearchDone from "@/pages/SearchDone";
import Document from "@/pages/Document/Document";

function App() {
  return (
    <Router>
      <AsideBar>
        <Routes>
          <Route path="/" element={<div>Главная страница</div>} />
          <Route path="/search-done" element={<SearchDone />} />
          <Route path="/document/:id" element={<Document />} />
          <Route path="*" element={<div>Страница не найдена</div>} />
        </Routes>
      </AsideBar>
    </Router>
  );
}

export default App;
