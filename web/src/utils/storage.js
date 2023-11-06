import Cookies from 'js-cookie';

export function StorageSetItem(key, val) {
    if (window['localStorage'] != null) {
        localStorage.setItem(key, val)
        return
    }
    // 兜底使用Cookies
    Cookies.set(key, val)
}

export function StorageGetItem(key) {
    if (window['localStorage'] != null) {
        return localStorage.getItem(key)
    }
    // 兜底使用Cookies
    return Cookies.get(key)
}