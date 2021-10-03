import { POSTRequestWithParams } from "utils/api";
import { ClosePopup } from "components/popup/popup";
import { Notify } from "components/app-notification/notification";

import styled from "styled-components";

const SToUp = styled.div`
    padding: 1rem;
    margin: 1rem;

    & .price {
        color: red;
        font-size: 1.3rem;
        background: yellow;
    }

    & .answer {
        display: flex;
        align-items: center;
        justify-content: space-evenly;

        & span {
            padding: .5rem 1rem;
            margin: 1rem;
            color: var(--onHoverColor);
            background: var(--blueColor);
            border-radius: 10px;
            cursor: pointer;
            transition: var(--transtionApp);

            &:nth-child(2) {
                background: red;
            }

            &:hover{
                background: var(--onHoverBG);
            }
        }
    }
`;

const toUp = async (id, type, code, cb) => {
    const res = await POSTRequestWithParams("/e/up", { 'id': id, 'type': type, 'code': code })
    if (res.err && res.err !== "ok") return Notify('fail', 'Не удалось поднять');
    cb()
    ClosePopup();
}

/**
 * 
 * @param type if cost will be relative by type: parsel or travel 
 * @param cb callback after click to up
 * @param id parsel/travel id 
 */
export default function ToUp({ cb, type, id }) {
    // const [price, setPrice] = useState();
    // const [isPayed, setPayed] = useState();
    // const code = useInput('');

    // const getPrice = useCallback(async () => {
    //     const res = await GetDataByCrieteries('price');
    //     if (res?.err !== "ok") return Notify('fail', "Ошибка. Попробуйте позднее") || setPrice(ExamplePrice)
    //     setPrice(res)
    // }, [])

    // useEffect(() => {
    //     if (!price) return getPrice()
    // }, [getPrice, price])


    // if (!price) return <div>Ошибка. Попробуйте позднее</div>
    return (
        <SToUp>
            <h2>Поднять Ваше объявление?</h2>
            <span>Поднимать можно только один раз в день</span>

            {/* <div className="price">
                Стоимость: {price?.cost} тг
            </div> */}

            <div className="answer">
                <span onClick={() => toUp(id, type, "", cb)}>Да</span>
                <span onClick={ClosePopup}>Нет</span>
            </div>

            {/* {
                isPayed &&
                <div>
                    <Input index="2" id="code" type="text" name="code" base={code.base} labelText="Введите 8-значный код:"
                        minLength="8" maxLength="8" placeholder="Mfa7sd45"
                    />

                    <SubmitBtn value="Поднять!" onClick={() => toUp(id, type, code, cb)} />
                </div>
            } */}
        </SToUp>
    )
}