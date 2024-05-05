import React from 'react';
import { Button } from 'react-native';

export const PressComponent = () => {
    return (
            <Button title="Button" onPress={() => alertText()} />
    );
}

const alertText = () => {
    const text: string = 'You pressed the button';
    alert(text);
    console.log("exit");
  }