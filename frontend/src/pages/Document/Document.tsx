import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getPageInfo } from "./getDocument";

export default function PageDetails() {
  const { id } = useParams<{ id: string }>();
  const [page, setPage] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    const fetchData = async () => {
      const res = await getPageInfo(id);
      if (res.ok && res.data) {
        setPage(res.data);
      }
      setLoading(false);
    };
    fetchData();
  }, [id]);

  if (loading) return <p>Загрузка...</p>;
  if (!page) return <p>Данные не найдены</p>;

  return (
    <div>
      <h1>{page.title}</h1>
      <p>{page.description}</p>
      <div dangerouslySetInnerHTML={{ __html: page.content }} />
    </div>
  );
}
