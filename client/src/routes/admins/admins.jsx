import { useCallback, useEffect, useRef, useState } from "react"
import { Redirect } from "react-router"

import { USER } from "constants/constants"
import { Fetching, GetDataByCrieteries } from "utils/api"
import { RandomKey } from "utils/content"
import { useInput } from "utils/form"
import { Notify } from "components/app-notification/notification"
import Input from "components/form-input/input"
import styled from "styled-components"

const SAdmin = styled.div`
    padding: 1rem;
    
    &.admins {
        margin-bottom: 7rem;
    }

    .btn {
        box-shadow: var(--boxShadow);
    }

    & > * {
        margin: 1rem;
    }

    .title {
        text-align: center;
    }

    .routes-items {
        max-height: 20vh;
        padding: 1rem;
        overflow: auto;

        .route-item {
            margin: 5px 0;
            padding: 5px;
            color: white;
            background: #192955;
            box-shadow: var(--boxShadow);
            cursor: pointer;
            transition: .5s;

            &:hover,
            &.active {
                background: white;
                color: #192955;
            }
        }
    }

    .active_route {
        border: 1px solid white;
        padding: 1rem;

        .active_title{
            text-align: center;

            b {
                color: red;
            }
        }
    }

    .method-items {
        display: flex;
        align-items: center;

        & span {
            padding: 1rem;
            margin: 1rem;
            border-radius: 10px;
            cursor: pointer;
            transition: .5s;

            &:hover,
            &.active{
                background: white;
                color: #192955;
            }
        }
    }

    .add-params {
        display: flex;
        align-items: center;

        & > * {
            margin: 1rem;
        }
    }

    .param-item {
        padding: .5rem 0;    
        display: flex;
        align-items: center;

        .param-content {
            display: flex;
            align-items: center;
            width: 80%;
        }

        .param-content > * {
            margin: 0 1rem;
            padding: 1rem;
            width: 50%;
            color: white;
            border-radius: 10px;
            background: #192955;
            box-shadow: var(--boxShadow);
        }
    }

    button.btn.btn-primary.search {
        margin: 1rem;
        width: 80%;
    }

    .results {
        color: white;

        pre {
            background: black;
        }
    }
`;

const GRoute = ({ active, route, description, methods, children = [], params = {}, setRoute }) => {
    return (
        <div className={`route-item ${active?.route === route ? "active" : ""}`} onClick={() => setRoute({ 'route': route, 'methods': methods, 'params': params ? params : undefined })}>
            <span>Путь <b> {route} </b> {description}</span>

            {children && children.map(r => <GRoute setRoute={setRoute} key={RandomKey()} {...r} />)}
        </div>
    )
}

const GParam = ({ k, v, removeParam }) => {
    return (
        <div className="param-item">
            <div className="param-content">
                <span>Ключ: {k}</span>
                <span>Значение: {v}</span>
            </div>

            {removeParam && <button className="btn btn-primary remove-param" onClick={() => removeParam(k)}>X</button>}
        </div>
    )
}

const GResultData = ({ data }) => {
    const rf = useRef(null);

    useEffect(() => {
        if (rf.current) rf.current.innerHTML = window.prettyPrintJson.toHtml(data)
    })

    return <pre ref={rf}></pre>
}

const addParam = (key, value, params, setParams, fields) => {
    setParams([...params, { 'k': key, 'v': value }])
    fields.forEach(f => f.resetField())
}

const removeParam = (key, params, setParams) => setParams([...params.filter(p => p.k !== key)])

export default function AdminPage() {
    const [routes, setRoutes] = useState()
    const [data, setData] = useState()

    const [route, setRoute] = useState()
    const [method, setMethod] = useState("GET")
    const [params, setParams] = useState([])

    const key = useInput("")
    const value = useInput("")
    const fields = [key, value];

    const getRoutes = useCallback(async () => {
        const resp = await GetDataByCrieteries("", {}, "GET")
        if (resp.err && resp.err !== "ok") return Notify('fail', "Не загрузились точки входа") || setRoutes(null);
        setRoutes(resp);
    }, [])

    const getData = useCallback(async () => {
        const data = new FormData();
        params.forEach(p => data.append(p.k, p.v))

        const resp = await Fetching("/api/" + route.route, data, method)
        if (resp.code !== 200) return Notify('fail', "Не загрузились данные") || setData(null);
        setData(resp.data);
    }, [route, params, method])

    useEffect(() => {
        if (routes === undefined) getRoutes()
    })

    if (!USER.isAdmin) return <Redirect to="/parsel" />
    if (!routes) return <div>Не загрузились точки входа</div>

    return (
        <SAdmin className="admins">
            <h1 className="title">Админка</h1>

            <div className="routes">
                <h2>Пути</h2>

                <div className="routes-items">
                    {
                        routes && routes.map(r => <GRoute active={route} setRoute={setRoute} key={RandomKey()} {...r} />)
                    }
                </div>
            </div>

            {
                route &&
                <div className="active_route">
                    <h3 className="active_title">Активный путь <b> {route.route} </b></h3>

                    <div className="methods">
                        <h3>Методы</h3>

                        <div className="method-items">
                            {route.methods.map(m => <span key={RandomKey()} className="btn btn-primary">{m}</span>)}
                        </div>
                    </div>

                    {
                        route?.params &&
                        <div className="params">
                            <h3>Параметры</h3>

                            <div className="params-items">
                                {route?.params.map(p => <GParam key={RandomKey()} {...p} />)}
                            </div>
                        </div>
                    }
                </div>
            }

            <div className="method">
                <h2>Метод</h2>

                <div className="method-items">
                    <span className={`btn btn-primary ${method === "GET" ? "active" : ""}`} onClick={() => setMethod("GET")}>GET</span>
                    <span className={`btn btn-primary ${method === "POST" ? "active" : ""}`} onClick={() => setMethod("POST")}>POST</span>
                </div>
            </div>

            <div className="params">
                <h2>Параметры</h2>

                <div className="add-params">
                    <Input type="text" labelText="Ключ" base={key.base} />
                    <Input type="text" labelText="Значение" base={value.base} />
                    <button className="btn btn-primary" onClick={() => addParam(key.base.value, value.base.value, params, setParams, fields)}>Добавить</button>
                </div>

                {params && params.map(p => <GParam key={RandomKey()} removeParam={(key) => removeParam(key, params, setParams)} {...p} />)}
            </div>

            <button className="btn btn-primary search" onClick={getData}>Поиск</button>

            {
                data &&
                <div className="results">
                    <h2>Результаты</h2>

                    {/* pretty print */}
                    <GResultData data={data} />
                </div>
            }

        </SAdmin>
    )
}