import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getPageInfo } from "./getDocument";
import { Button } from "@gravity-ui/uikit";
import { indexatePage, type V1DiagnosticInfoGetResponse } from "@/client";

export default function PageDetails() {
  const { id } = useParams<{ id: string }>();
  const [page, setPage] = useState<V1DiagnosticInfoGetResponse | undefined>(
    undefined
  );
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    const fetchData = async () => {
      const res = await getPageInfo(id);
      if (res.data) {
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
      <Button
        view="action"
        onClick={() => indexatePage({ body: { page_id: id! } })}
      >
        Проиндексировать
      </Button>
      <h1>{page.page.title}</h1>
      <p>{page.page.content}</p>
    </div>
  );
}
