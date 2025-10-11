import { search } from "@/client";
import styles from "./Search.module.scss";
import CardInformation from "@/components/CardInformation/CardInformation";
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function Search() {
  const [word, setWord] = useState("");
  const [results, setResults] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSearch = async () => {
    setLoading(true);
    const res = await search({ body: { query: word } });
    setResults(!res.error && res.data?.result_items ? res.data.result_items : []);
    setLoading(false);
  };

  const handleKeyPress = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") handleSearch();
  };

  const handleOpenPage = (pageId: string) => {
    navigate(`/document/${pageId}`);
  };

  return (
    <div className={styles.searchPage}>
      <input
        type="text"
        placeholder="Поиск..."
        value={word}
        onChange={(e) => setWord(e.target.value)}
        onKeyDown={handleKeyPress}
      />

      {loading && <p>Загрузка...</p>}

      {!loading && results.length > 0 && (
        <div className={styles.results}>
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
    </div>
  );
}
