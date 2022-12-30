import React, { ChangeEvent, useState } from 'react';
import { DeleteButton, Icon, IconButton, InlineField, InlineFieldRow, Input } from '@grafana/ui';

interface DerivedFiedProps {
    name: string
    value: string
    onChange: (value: string) => void;
    onDelete: () => void;
}
export const DerivedField: React.FC<DerivedFiedProps> = (props) => {
    const { name, onChange, onDelete, value } = props;

    const [state, setState] = useState({ curValue: value, });
    const { curValue } = state

    return (
        <InlineField label={name} transparent>
            <Input
                key={name}
                prefix={<Icon name='brackets-curly' />}
                suffix={<DeleteButton
                    closeOnConfirm
                    size="sm"
                    onConfirm={onDelete}
                />}
                placeholder='mustache template for field value'
                value={curValue}
                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                    setState({ ...state, curValue: event.target.value });
                }}
                onKeyDown={(e) => {
                    if (e.key === 'Enter') {
                        e.preventDefault();
                    }
                }}
                onKeyUp={(e) => {
                    e.preventDefault();
                    if (e.key === 'Enter') {
                        onChange(curValue);
                    }
                }}
                onBlur={(e) => {
                    onChange(curValue)
                }}
            />
        </InlineField>
    );
}

interface AddDerivedFieldsProps {
    onConfirm: (name: string) => void;
}
export const AddDerivedField: React.FC<AddDerivedFieldsProps> = (props) => {
    const { onConfirm } = props;

    const [state, setState] = useState({ name: "", });

    const { name } = state;
    const confirm = () => {
        if (name !== '') {
            onConfirm(name);
        }
        setState({ ...state, name: "" });
    }

    return (
        <InlineField>
            <Input
                placeholder='add field...'
                suffix={<IconButton
                    name='plus'
                    variant='primary'
                    onClick={confirm}
                />}
                value={name}
                onChange={(event: ChangeEvent<HTMLInputElement>) => {
                    setState({ ...state, name: event.target.value });
                }}
                onKeyDown={(e) => {
                    // onKeyDown is triggered before onKeyUp, triggering submit behaviour on Enter press if this component
                    // is used inside forms. Moving onKeyboardAdd callback here doesn't work since text input is not captured in onKeyDown
                    if (e.key === 'Enter') {
                        e.preventDefault();
                    }
                }}
                onKeyUp={(e) => {
                    e.preventDefault();
                    if (e.key === 'Enter') {
                        confirm();
                    }
                }}
            />
        </InlineField>
    );
}

export interface DerivedFieldData {
    name: string;
    value: string;
}

interface DerivedFieldsProps {
    fields: DerivedFieldData[];
    onChange: (fields: DerivedFieldData[]) => void;
}
export const DerivedFields: React.FC<DerivedFieldsProps> = (props) => {
    const { fields, onChange } = props;

    return (
        <InlineFieldRow>
            {fields.map((v, i) => (
                <DerivedField
                    key={v.name}
                    name={v.name}
                    value={v.value}
                    onChange={(v) => {
                        fields[i].value = v
                        onChange(fields)
                    }}
                    onDelete={() => {
                        fields.splice(i, 1)
                        onChange(fields)
                    }}
                />
            ))}
            <AddDerivedField
                onConfirm={(name) => {
                    onChange([...fields, { name, value: "" }])
                }}
            />
        </InlineFieldRow>
    );
}
