/**
 * 数据格式转换工具
 * 用于将JSON格式转换为multipart/form-data格式
 */

/**
 * 将JSON对象转换为FormData对象
 * @param data JSON对象
 * @returns FormData对象
 */
export const jsonToFormData = (data: Record<string, any>): FormData => {
  const formData = new FormData();

  const processData = (obj: any, prefix: string = '') => {
    if (obj === null || obj === undefined) {
      return;
    }

    if (typeof obj === 'object' && !(obj instanceof File) && !(obj instanceof Blob)) {
      Object.keys(obj).forEach((key) => {
        const fullKey = prefix ? `${prefix}[${key}]` : key;
        processData(obj[key], fullKey);
      });
    } else {
      formData.append(prefix, obj);
    }
  };

  processData(data);
  return formData;
};

/**
 * 将FormData对象转换为JSON对象
 * @param formData FormData对象
 * @returns JSON对象
 */
export const formDataToJson = (formData: FormData): any => {
  const data: any = {};

  formData.forEach((value, key) => {
    if (key.includes('[')) {
      const keys = key.split(/\[|\]/).filter(Boolean);
      let current = data;

      for (let i = 0; i < keys.length; i++) {
        const k = keys[i];
        if (i === keys.length - 1) {
          current[k] = value;
        } else {
          if (!current[k]) {
            current[k] = {};
          }
          current = current[k];
        }
      }
    } else {
      data[key] = value;
    }
  });

  return data;
};
