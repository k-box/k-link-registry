import store from "@/store";
import axios from "axios";

// true if page loaded via https
const ssl = (window.location.protocol === 'https:');

export function removePrefix(url) {
    if (url === '') url = '/';
    if (url[0] !== '/') url = '/' + url;
    return url;
}

export function fetch(url) {
    url = removePrefix(url);

    return new Promise((resolve, reject) => {
        let request = new window.XMLHttpRequest();
        request.open('GET', `${store.state.baseURL}/api/2.0${url}`, true);
        request.setRequestHeader('Authorization', `Bearer ${store.state.jwt}`);

        request.onload = () => {
            switch (request.status) {
                case 200:
                    resolve(JSON.parse(request.responseText));
                    break;
                default:
                    reject(new Error(request.status));
                    break;
            }
        };
        request.onerror = (error) => reject(error);
        request.send();
    });
}

export function remove(url) {
    url = removePrefix(url);

    return new Promise((resolve, reject) => {
        let request = new window.XMLHttpRequest();
        request.open('DELETE', `${store.state.baseURL}/api/2.0${url}`, true);
        request.setRequestHeader('Authorization', `Bearer ${store.state.jwt}`);

        request.onload = () => {
            if (request.status === 200) {
                resolve(request.responseText);
            } else {
                reject(request.responseText);
            }
        };

        request.onerror = (error) => reject(error);
        request.send();
    });
}

export function post(url, content = '', overwrite = false, onupload) {
    url = removePrefix(url);

    return new Promise((resolve, reject) => {
        let request = new window.XMLHttpRequest();
        request.open('POST', `${store.state.baseURL}/api/2.0${url}`, true);
        request.setRequestHeader('Authorization', `Bearer ${store.state.jwt}`);

        if (typeof onupload === 'function') {
            request.upload.onprogress = onupload;
        }

        request.onload = () => {
            if (request.status === 200) {
                resolve(request.responseText);
            } else if (request.status === 409) {
                reject(request.status);
            } else {
                reject(request.responseText);
            }
        };

        request.onerror = (error) => {
            reject(error);
        };
        request.send(content);
    });
}

// APPLICATIONS
export function getApplications() {
    return new Promise((resolve, reject) => {
        axios
            .get(`${store.state.baseURL}/api/2.0/applications`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function getApplication(id) {
    return new Promise((resolve, reject) => {
        axios
            .get(`${store.state.baseURL}/api/2.0/applications/${id}`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function newApplication(application) {
    return new Promise((resolve, reject) => {
        axios
            .post(`${store.state.baseURL}/api/2.0/applications`, application, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function updateApplication(application) {
    return new Promise((resolve, reject) => {
        axios
            .put(`${store.state.baseURL}/api/2.0/applications/${application.id}`, application, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function deleteApplication(id) {
    return new Promise((resolve, reject) => {
        axios
            .delete(`${store.state.baseURL}/api/2.0/applications/${id}`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve();
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

// K-Links
export function getKlinks() {
    return new Promise((resolve, reject) => {
        axios
            .get(`${store.state.baseURL}/api/2.0/klinks`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function getKlink(id) {
    return new Promise((resolve, reject) => {
        axios
            .get(`${store.state.baseURL}/api/2.0/klinks/${id}`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function newKlink(klink) {
    return new Promise((resolve, reject) => {
        axios
            .post(`${store.state.baseURL}/api/2.0/klinks`, klink, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function updateKlink(klink) {
    return new Promise((resolve, reject) => {
        axios
            .put(`${store.state.baseURL}/api/2.0/klinks/${klink.id}`, klink, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function deleteKlink(id) {
    return new Promise((resolve, reject) => {
        axios
            .delete(`${store.state.baseURL}/api/2.0/klinks/${id}`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve();
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

// REGISTRANTS
export function getRegistrants() {
    return new Promise((resolve, reject) => {
        axios
            .get(`${store.state.baseURL}/api/2.0/registrants`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function getRegistrant(id) {
    return new Promise((resolve, reject) => {
        axios
            .get(`${store.state.baseURL}/api/2.0/registrants/${id}`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function getEmailVerification(token) {
    return new Promise((resolve, reject) => {
        axios
            .get(`${store.state.baseURL}/api/2.0/auth/email-verification/${token}`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function performEmailVerification(emailVerification) {
    return new Promise((resolve, reject) => {
        axios
            .put(`${store.state.baseURL}/api/2.0/auth/email-verification/${emailVerification.token}`, emailVerification, {})
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function performPasswordReset(passwordReset) {
    return new Promise((resolve, reject) => {
        axios
            .put(`${store.state.baseURL}/api/2.0/auth/password-reset/${passwordReset.token}`, passwordReset, {})
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}


export function newRegistrant(registrant) {
    return new Promise((resolve, reject) => {
        axios
            .post(`${store.state.baseURL}/api/2.0/registrants`, registrant, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function updateRegistrant(registrant) {
    return new Promise((resolve, reject) => {
        axios
            .put(`${store.state.baseURL}/api/2.0/registrants/${registrant.id}`, registrant, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

export function deleteRegistrant(id) {
    return new Promise((resolve, reject) => {
        axios
            .delete(`${store.state.baseURL}/api/2.0/registrants/${id}`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve();
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

// PERMISSIONS
export function getPermissions() {
    return new Promise((resolve, reject) => {
        axios
            .get(`${store.state.baseURL}/api/2.0/permissions`, {
                headers: {
                    Authorization: `Bearer ${store.state.jwt}`
                }
            })
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}

// REGISTRATION
export function newRegistration(registration) {
    return new Promise((resolve, reject) => {
        axios
            .post(`${store.state.baseURL}/api/2.0/auth/registration`, registration, {})
            .then(response => {
                switch (response.status) {
                    case 200:
                        resolve(response.data);
                        break;
                    default:
                        reject(response.data.error);
                        break;
                }
            })
            .catch(e => {
                reject(e);
            });
    });
}