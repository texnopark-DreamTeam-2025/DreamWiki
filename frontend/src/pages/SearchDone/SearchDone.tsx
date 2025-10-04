import Search from "@/modules/Search/Search";

export default function PageSearch() {
  return (
    <div>
      <label>Страница поиска</label>
      <Search />
      <button>Поиск</button>
      <Search />
    </div>
  );
}
