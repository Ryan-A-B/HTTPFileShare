import React from 'react'

interface Props extends Omit<React.InputHTMLAttributes<HTMLInputElement>, "onChange"> {
    onChange: (value: string) => void
}

const TextInput: React.FunctionComponent<Props> = ({ onChange, ...props }) => {
    const onChangeInternal = React.useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
        onChange(event.target.value)
    }, [onChange])
    return (
        <input type="text" onChange={onChangeInternal} {...props} />
    )
}

export default TextInput