import { useCallback, useEffect, useState } from "react";

import { GetDataByCrieteries } from "utils/api";
import { RandomKey } from "utils/content";
import Input from "components/form-input/input";
import styled from "styled-components";

const SLalel = styled.label`
    display: flex;
    align-items: center;
    white-space: nowrap;

    & select {
        margin-left: 1rem;
        padding: .5rem;
        width: 100%;
        font-size: 1rem;
        color: var(--offHoverColor);
        background: none;
        border: none;
        border-radius: 5px;
        outline: none;
        border-bottom: 1px solid var(--onHoverColor);
        box-shadow: 4px 4px 3px 0 #00000029;
    }
`;

export default function PhoneField({ index, base, required }) {
    const [codes, setCodes] = useState();

    const getCodes = useCallback(async () => {
        const res = await GetDataByCrieteries('countryCodes');
        if (res.err && res?.err !== "ok") return setCodes(null);
        return setCodes(res)
    }, [])

    useEffect(() => {
        if (codes === undefined) return getCodes()
    }, [getCodes, codes])


    if (!codes) return null;
    return (
        <div>
            <Input index={index} id="phone" type="tel" name="phone" base={base} labelText="Телефон:"
                minLength="10" maxLength="15" placeholder="7777777777" required={required}
            />

            <SLalel>
                <span>Код страны:</span>
                <select name="countryCode">{codes?.map(({ code, country }) => <option key={RandomKey()} value={code}>{code} ({country})</option>)}</select>
            </SLalel>
        </div>
    )
}