import React from 'react';
import { create } from 'react-test-renderer';
import Main from '../../../pages/Main/Main';


describe("Main page", () => {
  it("renders correctly", () => {
    const tree = create(<Main />);
    expect(tree.toJSON()).toMatchSnapshot();
  });

  it("Genrates a random device name correctly", () => {

    const Main = require('../../../pages/Main/Main');

    // const mock = jest.spyOn(Main, "getRandomDeviceName");
    expect(Main.getRandomDeviceName()).toBeTruthy();
    expect(Main.getRandomDeviceName().length).not.toEqual(0);
    expect(typeof Main.getRandomDeviceName()).toBe("string");
  });

  it("Genrates a random Sensor correctly", () => {

    const Main = require('../../../pages/Main/Main');

    expect(Main.getRandomSensor()).toBeTruthy();
    expect(typeof Main.getRandomSensor()).toBe("object");

  });

  it("Genrates a random Value correctly", () => {

    const Main = require('../../../pages/Main/Main');

    expect(Main.getRandomValue(0, 1000)).toBeTruthy();
    expect(typeof Main.getRandomValue(-10, 1000)).toBe("number");

  });


});