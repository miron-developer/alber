const formDataToString = (data = new FormData()) => {
    let res = "";
    for (let [k, v] of data.entries())
        res += k + "=" + v + "&"
    return res.slice(0, -1)
}

// use fetching by both method
export const Fetching = async(action, data, method = "POST") => {
    if (action === undefined) return { err: "action undefined" };

    const fetchOption = { 'method': method };
    if (method === "GET") action += "?" + encodeURI(formDataToString(data));
    else fetchOption["body"] = data;

    return await fetch(action, fetchOption)
        .then(res => res.json())
        .catch(err => Object.assign({}, { 'err': "500: " + err }));
}

// convert from object to FormData
const prepareDataToFetch = (datas = {}) => {
    const data = new FormData();
    for (let [k, v] of Object.entries(datas)) data.append(k, v);
    return data;
}

// get data by criteries & type
export const GetDataByCrieteries = async(datatype, criteries = {}) => {
    const data = prepareDataToFetch(criteries);
    const res = await Fetching("/api/" + datatype, data, 'GET');
    if (res.code !== 200) return { 'err': res.err }
    return res.data;
}

// send post req to host with params
export const POSTRequestWithParams = async(to, params = {}) => {
    const data = prepareDataToFetch(params);
    return await Fetching(to, data);
}