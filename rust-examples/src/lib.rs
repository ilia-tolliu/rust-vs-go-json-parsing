use serde::{Deserialize, Serialize};

#[derive(Debug, PartialEq, Deserialize, Serialize)]
enum Fruit {
    #[serde(rename = "apple")]
    Apple,
    #[serde(rename = "orange")]
    Orange,
    #[serde(rename = "banana")]
    Banana,
}

#[derive(Deserialize, Serialize, Debug, PartialEq)]
struct JsonWithFruit {
    fruit: Fruit,
    owner: String,

    #[serde(skip_serializing_if = "Option::is_none")]
    description: Option<String>,
}

#[cfg(test)]
mod tests {
    use super::Fruit::Apple;
    use super::*;
    use anyhow::Result;
    use indoc::indoc;
    use serde_json::from_str;

    #[test]
    fn test_valid_json() -> Result<()> {
        let original_json = indoc! {r#"
            {
              "fruit": "apple",
              "owner": "John",
              "description": "a sweet one"
            }"#
        };

        let typed_value: JsonWithFruit = from_str(original_json)?;
        assert_eq!(
            typed_value,
            JsonWithFruit {
                fruit: Apple,
                owner: "John".to_string(),
                description: Some("a sweet one".to_string())
            }
        );

        let resulting_json = serde_json::to_string_pretty(&typed_value)?;
        assert_eq!(resulting_json, original_json);

        Ok(())
    }

    #[test]
    fn test_valid_json_without_optional_field() -> Result<()> {
        let original_json = indoc! {r#"
            {
              "fruit": "apple",
              "owner": "John"
            }"#
        };

        let typed_value: JsonWithFruit = from_str(original_json)?;
        assert_eq!(
            typed_value,
            JsonWithFruit {
                fruit: Apple,
                owner: "John".to_string(),
                description: None
            }
        );

        let resulting_json = serde_json::to_string_pretty(&typed_value)?;
        assert_eq!(resulting_json, original_json);

        Ok(())
    }

    #[test]
    fn test_invalid_json_required_enum_field_missing() {
        let original_json = indoc! {r#"
            {
              "owner": "John"
            }"#
        };

        let Err(err) = from_str::<JsonWithFruit>(original_json) else {
            panic!("didn't get error as expected");
        };

        assert!(err.to_string().contains("missing field `fruit`"));
    }

    #[test]
    fn test_invalid_json_required_string_field_missing() {
        let original_json = indoc! {r#"
            {
              "fruit": "apple"
            }"#
        };

        let Err(err) = from_str::<JsonWithFruit>(original_json) else {
            panic!("didn't get error as expected");
        };

        assert!(err.to_string().contains("missing field `owner`"));
    }

    #[test]
    fn test_invalid_json_wrong_enum_value() {
        let original_json = indoc! {r#"
            {
              "fruit": "appleWithTypo",
              "owner": "John"
            }"#
        };

        let Err(err) = from_str::<JsonWithFruit>(original_json) else {
            panic!("didn't get error as expected");
        };

        assert!(err.to_string().contains("unknown variant `appleWithTypo`"));
    }
}
