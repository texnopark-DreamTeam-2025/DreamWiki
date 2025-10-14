import { Button } from "@gravity-ui/uikit";
import { useToast } from "../../hooks/useToast";
import styles from "./HomePage.module.scss";

export const HomePage = () => {
  const { showSuccess, showError, showWarning, showInfo } = useToast();

  const handleSuccessToast = () => {
    showSuccess("Успех!", "Операция выполнена успешно");
  };

  const handleErrorToast = () => {
    showError("Ошибка!", "Что-то пошло не так");
  };

  const handleWarningToast = () => {
    showWarning("Предупреждение!", "Обратите внимание на это");
  };

  const handleInfoToast = () => {
    showInfo("Информация", "Полезная информация для пользователя");
  };

  return (
    <div className={styles.homePage}>
      <h1>Главная страница</h1>
      <p>Добро пожаловать в DreamWiki!</p>

      <div className={styles.toastExamples}>
        <h3>Примеры уведомлений (Toaster):</h3>
        <div className={styles.buttonGroup}>
          <Button view="action" onClick={handleSuccessToast}>
            Успех
          </Button>
          <Button view="outlined-danger" onClick={handleErrorToast}>
            Ошибка
          </Button>
          <Button view="outlined-warning" onClick={handleWarningToast}>
            Предупреждение
          </Button>
          <Button view="outlined-info" onClick={handleInfoToast}>
            Информация
          </Button>
        </div>
      </div>
    </div>
  );
};
