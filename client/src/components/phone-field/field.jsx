import { useCallback, useEffect, useState } from "react";

import { GetDataByCrieteries } from "utils/api";
import Input from "components/form-input/input";
import Select from "components/form-select/select";

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


    if (!codes) return <div></div>;
    return (
        <div>
            <Input index={index} id="phone" type="number" name="phone" base={base} labelText="Телефон:"
                minLength="10" maxLength="15" placeholder="7777777777" required={required}
            />

            <Select text="Код страны:" name="countryCode" required={required} options={{
                data: codes,
                value: "code",
                makeText: ({code, country}) => `${code} (${country})`
            }} />
        </div>
    )
}