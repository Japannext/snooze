'''Tests for model related utils'''

import pytest
from pydantic import BaseModel, ValidationError

from snooze.utils.condition import *
from snooze.utils.model import *

class TestPartial:
    def test_partial_init(self):
        new_model = Partial[Rule]
        print(Rule.__fields__)
        print(new_model.__fields__)
        for name, field in new_model.__fields__.items():
            assert field.required == False
            assert field.default == None
        data = new_model()
        assert data.uid == None
        assert data.name == None
        assert data.comment == None
        assert data.condition == None
        assert isinstance(data, BaseModel)
        assert isinstance(data, Partial[Rule])

    def test_partial_with_data_simple(self):
        new_model = Partial[Rule]

        data = new_model(comment='My new comment')
        assert data.uid == None
        assert data.name == None
        assert data.comment == 'My new comment'
        assert data.condition == None

    def test_partial_with_data_complex(self):
        new_model = Partial[Rule]

        data = new_model(condition=AlwaysTrue())
        assert data.uid == None
        assert data.name == None
        assert data.comment == None
        assert data.condition.type == 'ALWAYS_TRUE'

    def test_partial_with_data_2_uses(self):
        data1 = Partial[Rule](comment='my comment')
        assert data1.uid == None
        assert data1.name == None
        assert data1.comment == 'my comment'
        assert data1.condition == None

        data2 = Partial[Rule](condition=AlwaysTrue())
        assert data2.uid == None
        assert data2.name == None
        assert data2.comment == None
        assert data2.condition.type == 'ALWAYS_TRUE'

    def test_partial_validation_failure(self):
        new_model = Partial[Rule]
        with pytest.raises(ValidationError):
            data = new_model(comment={'type': 'unknown condition type', 'params': 'blah'})
