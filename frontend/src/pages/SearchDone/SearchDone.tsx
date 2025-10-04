import { useEffect, useState } from "react";
import styles from "./SearchDone.module.scss";
import {
  DefaultApi,
  createConfiguration,
  ServerConfiguration,
} from "../../codegen";
import type { SearchResponse } from "../../codegen";

export default function SearchDone() {
  const [data, setData] = useState<SearchResponse | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const api = new DefaultApi(
      createConfiguration({
        baseServer: new ServerConfiguration("http://zhugeo.ru", {}),
      })
    );

    const searchData = async () => {
      setLoading(true);
      try {
        const response = await api.search({ query: "example" });
        setData(response);
      } catch (error) {
        console.error("API Error:", error);
      } finally {
        setLoading(false);
      }
    };

    searchData();
  }, []);

  return (
    <div className={styles.container}>
      <h1>Search Done Page</h1>
      {loading && <p>Загрузка...</p>}
      {data && <pre>{JSON.stringify(data, null, 2)}</pre>}
      <button>Search</button>
    </div>
  );
}
