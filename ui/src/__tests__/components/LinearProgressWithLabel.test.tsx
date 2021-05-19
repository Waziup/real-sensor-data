import React from 'react';
import LinearProgressWithLabel from '../../components/LinearProgressWithLabel';
import { create } from 'react-test-renderer'

describe('snapshot test', () => {
    test('testing LinearProgressWithLabel', () => {
        let tree = create(<LinearProgressWithLabel value={50} />)
        expect(tree.toJSON()).toMatchSnapshot();
    })
})