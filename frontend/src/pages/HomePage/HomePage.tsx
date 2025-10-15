import Search from "@/modules/Search/Search";
import styles from "./HomePage.module.scss";

export const HomePage = () => {
  return (
    <div className={styles.homePage}>
      <div className={styles.welcomeSection}>
        <h1>DreamWiki</h1>
        <p>База знаний с интеллектуальным поиском</p>
      </div>

      <div className={styles.searchSection}>
        <Search />
      </div>
    </div>
  );
};
