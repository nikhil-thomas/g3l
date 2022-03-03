package g3l_maps

type Dict map[string]string

var (
    ErrNotFound         = DictionaryErr("could not find the word you were looking for")
    ErrWordExists       = DictionaryErr("cannot add word because it already exists")
    ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}

func (d Dict) Search(key string) (string, error) {
    val, ok := d[key]
    if !ok {
        return "", ErrNotFound
    }
    return val, nil
}

func (d Dict) Add(key, val string) error {
    _, err := d.Search(key)
    switch err {
    case ErrNotFound:
        d[key] = val
    case nil:
        return ErrWordExists
    default:
        return err
    }
    return nil
}

func (d Dict) Update(key, val string) error {
    _, err := d.Search(key)
    switch err {
    case ErrNotFound:
        return ErrWordDoesNotExist
    case nil:
        d[key] = val
    default:
        return err

    }
    return nil
}

func (d Dict) Delete(key string) error {
    _, err := d.Search(key)
    switch err {
    case nil:
        delete(d, key)
        return nil
    case ErrNotFound:
        return ErrWordDoesNotExist
    default:
        return err
    }
    return nil
}
