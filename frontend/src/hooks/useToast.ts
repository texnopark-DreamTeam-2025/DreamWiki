import { toaster } from "@gravity-ui/uikit/toaster-singleton";
import type { ToastProps } from "@gravity-ui/uikit";

export const useToast = () => {
  const showToast = (props: Omit<ToastProps, "name"> & { name?: string }) => {
    const toastName = props.name || `toast-${Date.now()}`;
    toaster.add({
      ...props,
      name: toastName,
    });
    return toastName;
  };

  const showSuccess = (title: string, content?: React.ReactNode) => {
    return showToast({
      title,
      content,
      theme: "success",
      autoHiding: 5000,
    });
  };

  const showError = (title: string, content?: React.ReactNode) => {
    return showToast({
      title,
      content,
      theme: "danger",
      autoHiding: 7000,
    });
  };

  const showWarning = (title: string, content?: React.ReactNode) => {
    return showToast({
      title,
      content,
      theme: "warning",
      autoHiding: 6000,
    });
  };

  const showInfo = (title: string, content?: React.ReactNode) => {
    return showToast({
      title,
      content,
      theme: "info",
      autoHiding: 5000,
    });
  };

  const removeToast = (name: string) => {
    toaster.remove(name);
  };

  const removeAllToasts = () => {
    toaster.removeAll();
  };

  return {
    showToast,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    removeToast,
    removeAllToasts,
  };
};
