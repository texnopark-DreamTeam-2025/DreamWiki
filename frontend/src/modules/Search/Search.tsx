import { search } from "@/client";
import type { SearchResultItem } from "@/client";
import { useToast } from "@/hooks/useToast";
import styles from "./Search.module.scss";
import CardInformation from "@/components/CardInformation/CardInformation";
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function Search() {
  const [word, setWord] = useState("");
  const [results, setResults] = useState<SearchResultItem[]>([]);
  const [loading, setLoading] = useState(false);
  const { showSuccess, showError } = useToast();
  const navigate = useNavigate();

  const handleSearch = async () => {
    if (!word.trim()) {
      showError("Ошибка", "Введите поисковый запрос");
      return;
    }

    setLoading(true);
    try {
      const res = await search({ body: { query: word.trim() } });
      
      if (res.error) {
        console.error("Ошибка поиска:", res.error);
        showError("Ошибка поиска", "Не удалось выполнить поиск. Попробуйте еще раз.");
        setResults([]);
        return;
      }

      const searchResults = res.data?.result_items || [];
      setResults(searchResults);
      
      if (searchResults.length === 0) {
        showSuccess("Поиск завершен", "По вашему запросу ничего не найдено");
      } else {
        showSuccess("Поиск завершен", `Найдено результатов: ${searchResults.length}`);
      }
    } catch (error) {
      console.error("Ошибка при выполнении поиска:", error);
      showError("Ошибка", "Произошла ошибка при выполнении поиска");
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") handleSearch();
  };

  const handleOpenPage = (pageId: string) => {
    navigate(`/document/${pageId}`);
  };

  return (
    <div className={styles.searchPage}>
      <div className={styles.searchContainer}>
        <input
          type="text"
          placeholder="Введите поисковый запрос на естественном языке..."
          value={word}
          onChange={(e) => setWord(e.target.value)}
          onKeyDown={handleKeyPress}
          className={styles.searchInput}
          disabled={loading}
        />
        <button 
          onClick={handleSearch}
          disabled={loading || !word.trim()}
          className={styles.searchButton}
        >
          {loading ? "Поиск..." : "Найти"}
        </button>
      </div>

      {loading && (
        <div className={styles.loading}>
          <p>Выполняется поиск...</p>
        </div>
      )}

      {!loading && results.length > 0 && (
        <div className={styles.results}>
          <div className={styles.resultsHeader}>
            <h3>Результаты поиска ({results.length})</h3>
          </div>
          {results.map((item) => (
            <div
              key={item.page_id}
              onClick={() => handleOpenPage(item.page_id)}
              className={styles.resultItem}
            >
              <CardInformation
                title={item.title}
                description={item.description}
              />
            </div>
          ))}
        </div>
      )}

      {!loading && word && results.length === 0 && (
        <div className={styles.noResults}>
          <p>По запросу "{word}" ничего не найдено</p>
          <p>Попробуйте изменить запрос или использовать другие ключевые слова</p>
        </div>
      )}
    </div>
  );
}
